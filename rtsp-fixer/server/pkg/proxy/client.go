package proxy

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
	"github.com/bluenviron/gortsplib/v4/pkg/description"
	"github.com/bluenviron/gortsplib/v4/pkg/format"
	"github.com/bluenviron/gortsplib/v4/pkg/format/rtph264"
	"github.com/bluenviron/gortsplib/v4/pkg/headers"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
)

const connTimeout = 4 * time.Second

type thumbnailRecorder struct {
	h264 *rtph264.Decoder
}

type thumbnailStream struct {
	media       *description.Media
	originalErr error
}

type client struct {
	srv               *server
	clt               gortsplib.Client
	path              string
	config            *StreamConfig
	stream            *gortsplib.ServerStream
	thumbnailRecorder thumbnailRecorder
	thumbnailStream   *thumbnailStream
	ctx               context.Context
	cancel            context.CancelFunc
	closeOnce         sync.Once
}

func (me *server) createClientFor(path string) (*client, error) {
	stream, ok := me.streams[path]
	if !ok {
		return nil, fmt.Errorf("stream %q not found", path)
	}

	ctx, cancel := context.WithCancel(context.Background())
	clt := &client{
		srv:    me,
		ctx:    ctx,
		cancel: cancel,
		clt: gortsplib.Client{
			Scheme:       stream.URL.Scheme,
			Host:         stream.URL.Host,
			UserAgent:    userAgent,
			ReadTimeout:  connTimeout,
			WriteTimeout: connTimeout,
		},
		path:   path,
		config: stream,
	}
	clt.clt.OnRequest = clt.onRequest
	clt.clt.OnResponse = clt.onResponse
	clt.clt.OnTransportSwitch = clt.onTransportSwitch
	clt.clt.OnPacketsLost = clt.onPacketsLost
	clt.clt.OnDecodeError = clt.onDecodeError
	err := clt.clt.Start2()
	if err != nil {
		return nil, err
	}

	return clt, nil
}

func (me *client) Close() error {
	me.closeOnce.Do(func() {
		me.cancel()
		if me.stream != nil {
			me.stream.Close()
		}
		me.clt.Close()
	})
	return nil
}

func (me *client) Describe(url *base.URL) (*description.Session, *base.Response, error) {
	me.srv.thumbsMutex.RLock()
	thumb, ok := me.srv.thumbs[me.path]
	var connErr error
	if ok {
		connErr = thumb.connErr
	}
	me.srv.thumbsMutex.RUnlock()

	if connErr != nil {
		return me.makeThumbnailStream(connErr)
	}
	session, resp, err := me.clt.Describe(url)
	if err != nil {
		me.srv.updateConnErr(me.path, err)
		return me.makeThumbnailStream(err)
	}
	return session, resp, err
}

func (me *client) Setup(url *base.URL, media *description.Media, rtpPort, rtcpPort int) (*base.Response, error) {
	if me.isThumbnail() {
		return &base.Response{StatusCode: base.StatusOK}, nil
	}
	return me.clt.Setup(url, media, rtpPort, rtcpPort)
}

func (me *client) Play(rng *headers.Range) (*base.Response, error) {
	logger := me.srv.logger.WithField("client", me.path)

	if me.isThumbnail() {
		return me.playThumbnailStream()
	}

	me.clt.OnPacketRTCPAny(func(m *description.Media, p rtcp.Packet) {
		err := me.stream.WritePacketRTCP(m, p)
		if err != nil {
			logger.WithError(err).Errorf("error writing RTCP packet from client")
		}
	})
	recordThumbnail := true
	me.clt.OnPacketRTPAny(func(m *description.Media, f format.Format, p *rtp.Packet) {
		err := me.stream.WritePacketRTP(m, p)
		if err != nil {
			logger.WithError(err).Errorf("error writing RTP packet from client")
		}

		if !recordThumbnail {
			return
		}

		record, err := me.recordThumbnail(f, p)
		if err != nil {
			logger.WithError(err).Errorf("error recording thumbnail")
		}
		if !record {
			recordThumbnail = false
			logger.Warnf("skipping thumbnail recording for this stream")
		}
	})

	return me.clt.Play(rng)
}

func (me *client) Pause() (*base.Response, error) {
	if me.isThumbnail() {
		return errNotImplemented, nil
	}
	return me.clt.Pause()
}

func (me *client) Record() (*base.Response, error) {
	if me.isThumbnail() {
		return errNotImplemented, nil
	}
	return me.clt.Record()
}

func (me *client) Announce(url *base.URL, desc *description.Session) (*base.Response, error) {
	if me.isThumbnail() {
		return errNotImplemented, nil
	}
	return me.clt.Announce(url, desc)
}

func (me *client) onRequest(req *base.Request) {
	me.srv.logger.Debugf("proxying request to client: %+v", req)
}

func (me *client) onResponse(resp *base.Response) {
	if me.config.FixForceTCPInTransport {
		transport, found := resp.Header["Transport"]
		me.srv.logger.Debugf("trying to fix transport header: %+v (found: %v)", transport, found)
		if found && len(transport) == 1 && strings.HasPrefix(transport[0], "RTP/AVP;") {
			me.srv.logger.Infof("fixing transport header")
			resp.Header["Transport"] = base.HeaderValue{strings.Replace(transport[0], "RTP/AVP;", "RTP/AVP/TCP;", 1)}
		}
	}
	me.srv.logger.Debugf("received response from client: %+v", resp)
}

func (me *client) onTransportSwitch(err error) {
	me.srv.logger.WithError(err).Warningf("transport switched")
}

func (me *client) onPacketsLost(packetsLost uint64) {
	me.srv.logger.Errorf("lost %d client packets", packetsLost)
}

func (me *client) onDecodeError(err error) {
	me.srv.logger.WithError(err).Errorf("client decode error")
}

func getOptionalClient(ctx *gortsplib.ServerConn) (*client, error) {
	if ctx.UserData() == nil {
		return nil, fmt.Errorf("internal error: client not found")
	}
	clt, cast := ctx.UserData().(*client)
	if !cast {
		return nil, fmt.Errorf("internal error: invalid client type")
	}
	return clt, nil
}

func getClient(ctx *gortsplib.ServerConn) *client {
	clt, err := getOptionalClient(ctx)
	if err != nil {
		panic(err)
	}
	return clt
}
