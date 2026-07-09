package proxy

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bluenviron/gortsplib/v4/pkg/base"
	"github.com/bluenviron/gortsplib/v4/pkg/description"
	"github.com/bluenviron/gortsplib/v4/pkg/format"
	"github.com/pion/rtp"
	xdraw "golang.org/x/image/draw"
	xfont "golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	thumbnailInterval       = 5 * 60 * time.Second
	upstreamMonitorInterval = 10 * time.Second
)

var jpegOptions = &jpeg.Options{Quality: 85}

type thumbnailData struct {
	lastCapture   time.Time
	thumbnailPath string
	thumbnail     image.Image
	connErr       error
}

func (me *server) shouldRecordThumbnail(path string) bool {
	me.thumbsMutex.Lock()
	defer me.thumbsMutex.Unlock()

	thumb, ok := me.thumbs[path]
	if !ok {
		return false
	}
	if time.Since(thumb.lastCapture) < thumbnailInterval {
		return false
	}
	thumb.lastCapture = time.Now()
	return true
}

func thumbnailPath(r *http.Request) string {
	return fmt.Sprintf("/%s", strings.TrimPrefix(r.URL.Path, "/"))
}

func (me *server) resetThumbnail(w http.ResponseWriter, r *http.Request) {
	path := thumbnailPath(r)

	me.thumbsMutex.Lock()
	thumb, ok := me.thumbs[path]
	if ok {
		thumb.lastCapture = time.Time{}
		me.thumbs[path] = thumb
	}
	me.thumbsMutex.Unlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	me.logger.Infof("thumbnail timeout reset for %q", path)
	w.WriteHeader(http.StatusNoContent)
}

func (me *server) serveThumbnail(w http.ResponseWriter, r *http.Request) {
	path := thumbnailPath(r)

	me.thumbsMutex.RLock()
	thumb, ok := me.thumbs[path]
	var img image.Image
	var connErr error
	if ok {
		img = thumb.thumbnail
		connErr = thumb.connErr
	}
	me.thumbsMutex.RUnlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	me.logger.Debugf("serving thumbnail for %q", path)

	if img == nil {
		me.logger.Warnf("thumbnail for %q not found, using black placeholder", path)
		img = getFallbackThumbnail()
	}
	if connErr != nil {
		img = overlayError(img, connErr)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	err := jpeg.Encode(w, img, jpegOptions)
	if err != nil {
		me.logger.WithError(err).Errorf("failed to encode thumbnail for %q", path)
	}
}

func (me *server) getThumbnail(path string) (image.Image, error) {
	me.thumbsMutex.RLock()
	defer me.thumbsMutex.RUnlock()

	thumb, ok := me.thumbs[path]
	if !ok {
		return nil, fmt.Errorf("thumbnail for %q not found", path)
	}

	return thumb.thumbnail, nil
}

func (me *server) saveThumbnailH264(path string, nalUnits [][]byte) {
	if !me.shouldRecordThumbnail(path) {
		return
	}
	me.logger.Infof("saving thumbnail for %s...", path)
	go func() {
		err := me.saveThumbnailInternal(path, func() (image.Image, error) {
			return h264ToJPEG(nalUnits)
		})
		if err != nil {
			me.logger.WithError(err).Errorf("failed to save thumbnail for %s", path)
		}
	}()
}

func (me *server) saveThumbnailInternal(path string, getImage func() (image.Image, error)) error {
	me.thumbsMutex.RLock()
	thumb, ok := me.thumbs[path]
	me.thumbsMutex.RUnlock()
	if !ok {
		return fmt.Errorf("thumbnail for %q not found", path)
	}

	img, err := getImage()
	if err != nil {
		return fmt.Errorf("failed to get image for thumbnail: %w", err)
	}

	tmp := thumb.thumbnailPath + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("failed to create thumbnail file: %w", err)
	}
	defer f.Close()
	err = jpeg.Encode(f, img, jpegOptions)
	if err != nil {
		return fmt.Errorf("failed to encode thumbnail: %w", err)
	}
	err = f.Close()
	if err != nil {
		return fmt.Errorf("failed to close thumbnail file: %w", err)
	}
	err = os.Rename(tmp, thumb.thumbnailPath)
	if err != nil {
		return fmt.Errorf("failed to rename thumbnail file: %w", err)
	}

	me.logger.Infof("thumbnail for %s saved to %s", path, thumb.thumbnailPath)
	me.thumbsMutex.Lock()
	thumb.thumbnail = img
	thumb.lastCapture = time.Now()
	me.thumbsMutex.Unlock()
	return nil
}

// client

func (me *client) makeThumbnailStream(originalErr error) (*description.Session, *base.Response, error) {
	me.srv.logger.WithError(originalErr).Warnf("stream %q appears to be offline, creating thumbnail stream", me.path)
	forma := &format.MJPEG{}
	media := &description.Media{
		Type:    description.MediaTypeVideo,
		Formats: []format.Format{forma},
	}
	desc := &description.Session{
		Medias: []*description.Media{media},
	}
	me.thumbnailStream = &thumbnailStream{
		originalErr: originalErr,
		media:       media,
	}
	return desc, &base.Response{StatusCode: base.StatusOK}, nil
}

func (me *client) playThumbnailStream() (*base.Response, error) {
	img, err := me.srv.getThumbnail(me.path)
	if img == nil {
		me.srv.logger.WithError(err).Warnf("no thumbnail for %q, using black placeholder", me.path)
		img = getFallbackThumbnail()
	}
	img = overlayError(img, me.thumbnailStream.originalErr)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, jpegOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to encode JPEG: %w", err)
	}

	enc, err := (&format.MJPEG{}).CreateEncoder()
	if err != nil {
		return nil, fmt.Errorf("failed to create encoder: %w", err)
	}

	jpegBytes := buf.Bytes()
	ch := me.srv.pubSub.subscribe(me.path)
	go func() {
		defer me.srv.pubSub.unsubscribe(me.path, ch)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		frames := uint32(0)

		for {
			select {
			case <-me.ctx.Done():
				return
			case connErr := <-ch:
				if connErr == nil {
					me.srv.logger.Infof("stream %q came back online, closing thumbnail stream", me.path)
					me.Close()
					return
				}
			case <-ticker.C:
				pkts, err := enc.Encode(jpegBytes)
				if err != nil {
					me.srv.logger.WithError(err).Errorf("thumbnail: failed to encode MJPEG frame")
					return
				}
				for _, pkt := range pkts {
					pkt.Timestamp = frames * 90000
					me.srv.logger.Debugf("thumbnail: sending RTP packet with seq=%d ts=%d marker=%v", pkt.SequenceNumber, pkt.Timestamp, pkt.Marker)
					err := me.stream.WritePacketRTP(me.thumbnailStream.media, pkt)
					if err != nil {
						me.srv.logger.WithError(err).Errorf("thumbnail: failed to write RTP packet")
						return
					}
				}
				frames += 1
			}
		}
	}()

	return &base.Response{StatusCode: base.StatusOK}, nil
}

func (me *client) isThumbnail() bool {
	return me.thumbnailStream != nil
}

func (me *client) recordThumbnail(f format.Format, p *rtp.Packet) (bool, error) {
	switch v := f.(type) {
	case *format.H264:
		return true, me.recordH264Thumbnail(v, p)
	case *format.G711:
		// ignore audio formats
		return true, nil
	default:
		return false, fmt.Errorf("unsupported format type for thumbnail: %T", f)
	}
}

func (me *client) recordH264Thumbnail(f *format.H264, p *rtp.Packet) error {
	if me.thumbnailRecorder.h264 == nil { // lazy
		d, err := f.CreateDecoder()
		if err != nil {
			return fmt.Errorf("failed to create H264 depacketizer: %w", err)
		}
		me.thumbnailRecorder.h264 = d
	}

	nalUnits, err := me.thumbnailRecorder.h264.Decode(p)
	if err != nil {
		// "need more packets" is normal for FU-A fragments, not a real error
		return nil
	}

	buffer := make([][]byte, 0)
loop:
	for _, nalu := range nalUnits {
		if len(nalu) == 0 {
			continue
		}
		naluType := nalu[0] & 0x1F
		switch naluType {
		case 7, 8: // SPS/PPS: keep the latest
			buffer = append(buffer, nalu)
		case 5: // IDR: we have everything we need, write and reset
			buffer = append(buffer, nalu)
			me.srv.saveThumbnailH264(me.path, buffer)
			break loop
		}
	}

	return nil
}

// helpers

func getThumbnailPath(base, path string) string {
	hashPathB := md5.Sum([]byte(path))
	hashPath := hex.EncodeToString(hashPathB[:])
	filePath := filepath.Join(base, fmt.Sprintf("%s.jpg", hashPath))
	return filePath
}

func readThumbnail(inputFile string) (image.Image, error) {
	f, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open thumbnail file: %w", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	return img, err
}

func getFallbackThumbnail() image.Image {
	return image.NewRGBA(image.Rect(0, 0, 640, 480))
}

func overlayError(img image.Image, err error) image.Image {
	const scale = 3
	const padding = 10

	face := basicfont.Face7x13
	charW := face.Advance
	charH := face.Height

	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, img, bounds.Min, draw.Src)

	// wrap to fit image width at scaled size
	msg := err.Error()
	maxChars := max(bounds.Dx()/(charW*scale), 1)
	var lines []string
	for len(msg) > maxChars {
		lines = append(lines, msg[:maxChars])
		msg = msg[maxChars:]
	}
	lines = append(lines, msg)

	maxLen := 0
	for _, l := range lines {
		if len(l) > maxLen {
			maxLen = len(l)
		}
	}

	// render text at 1x onto offscreen
	offscreen := image.NewRGBA(image.Rect(0, 0, maxLen*charW, len(lines)*charH))
	d := &xfont.Drawer{
		Dst:  offscreen,
		Src:  image.NewUniform(color.RGBA{255, 80, 80, 255}),
		Face: face,
	}
	for i, line := range lines {
		d.Dot = fixed.P(0, (i+1)*charH-face.Descent)
		d.DrawString(line)
	}

	// scale up and center horizontally, anchored to bottom
	scaledW := maxLen * charW * scale
	scaledH := len(lines) * charH * scale
	x := max(bounds.Min.X+(bounds.Dx()-scaledW)/2, bounds.Min.X)
	y := bounds.Max.Y - scaledH - padding

	bgRect := image.Rect(bounds.Min.X, y-padding, bounds.Max.X, bounds.Max.Y)
	draw.Draw(dst, bgRect, image.NewUniform(color.RGBA{0, 0, 0, 180}), image.Point{}, draw.Over)

	xdraw.NearestNeighbor.Scale(dst, image.Rect(x, y, x+scaledW, y+scaledH), offscreen, offscreen.Bounds(), xdraw.Over, nil)

	return dst
}
