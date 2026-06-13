package proxy

import (
	"bytes"
	"fmt"
	"image"

	"github.com/y9o/go-openh264"
)

func init() {
	err := openh264.Open(libOpenH264)
	if err != nil {
		panic(fmt.Sprintf("failed to load openh264: %v", err))
	}
}

var annexBStartCode = []byte{0, 0, 0, 1}

func h264ToJPEG(nalUnits [][]byte) (image.Image, error) {
	// build a single Annex B buffer from all NAL units
	var buf []byte
	for _, nalu := range nalUnits {
		buf = append(buf, annexBStartCode...)
		buf = append(buf, nalu...)
	}

	var ppdec *openh264.ISVCDecoder
	if ret := openh264.WelsCreateDecoder(&ppdec); ret != 0 || ppdec == nil {
		return nil, fmt.Errorf("failed to create decoder: %d", ret)
	}
	defer openh264.WelsDestroyDecoder(ppdec)

	sDecParam := openh264.SDecodingParam{}
	sDecParam.EEcActiveIdc = openh264.ERROR_CON_SLICE_MV_COPY_CROSS_IDR_FREEZE_RES_CHANGE
	if r := ppdec.Initialize(&sDecParam); r != 0 {
		return nil, fmt.Errorf("failed to initialize decoder: %d", r)
	}
	defer ppdec.Uninitialize()

	var img *image.YCbCr
	for len(buf) > 4 {
		// find next start code to delimit this NAL unit
		pos := bytes.Index(buf[4:], annexBStartCode)
		length := len(buf)
		if pos != -1 {
			length = pos + 4
		}

		var pDst [3][]byte
		var sDstBufInfo openh264.SBufferInfo
		if r := ppdec.DecodeFrameNoDelay(buf[:length], length, &pDst, &sDstBufInfo); r != 0 {
			buf = buf[length:]
			continue
		}
		if pDst[0] != nil {
			sys := sDstBufInfo.UsrData_sSystemBuffer()
			img = &image.YCbCr{
				Y:              pDst[0],
				Cb:             pDst[1],
				Cr:             pDst[2],
				YStride:        int(sys.IStride[0]),
				CStride:        int(sys.IStride[1]),
				SubsampleRatio: image.YCbCrSubsampleRatio420,
				Rect:           image.Rect(0, 0, int(sys.IWidth), int(sys.IHeight)),
			}
		}
		if pos == -1 {
			break
		}
		buf = buf[pos+4:]
	}

	if img == nil {
		return nil, fmt.Errorf("no frame decoded")
	}
	return img, nil
}
