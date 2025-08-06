package proxy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
	"github.com/bluenviron/gortsplib/v4/pkg/description"
	"github.com/bluenviron/gortsplib/v4/pkg/format"
	"github.com/bluenviron/gortsplib/v4/pkg/headers"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
)

var (
	errInternal = &base.Response{
		StatusCode:    base.StatusInternalServerError,
		StatusMessage: "Internal error",
	}
	errGatewayTimeout = &base.Response{
		StatusCode:    base.StatusGatewayTimeout,
		StatusMessage: "Gateway error",
	}
	errNotImplemented = &base.Response{
		StatusCode: base.StatusNotImplemented,
	}
)

func wrapError(rsp *base.Response, err error) *base.Response {
	reply := *rsp
	reply.Body = []byte(err.Error())
	return &reply
}

func (me *server) OnDescribe(ctx *gortsplib.ServerHandlerOnDescribeCtx) (*base.Response, *gortsplib.ServerStream, error) {
	me.logger.Debugf("DESCRIBE request %+v", ctx)
	clt := getClient(ctx.Conn)

	url := clt.config.URL
	me.logger.Debugf("Proxying DESCRIBE request to %s", url.String())
	desc, res, err := clt.Describe((*base.URL)(&url))
	if err != nil {
		me.logger.WithError(err).Errorf("error proxying DESCRIBE request")
		return wrapError(errGatewayTimeout, err), nil, nil
	}

	me.logger.Debugf("Got SDP %+v", desc)
	for _, media := range desc.Medias {
		me.logger.Debugf("Got SDP Media %+v", *media)
	}

	if clt.stream == nil {
		clt.stream = &gortsplib.ServerStream{
			Server: me.rtsp,
			Desc:   desc,
		}
		err = clt.stream.Initialize()
		if err != nil {
			me.logger.WithError(err).Errorf("error initializing stream")
			return wrapError(errGatewayTimeout, err), nil, nil
		}
	}

	return res, clt.stream, nil
}

func (me *server) OnSetup(ctx *gortsplib.ServerHandlerOnSetupCtx) (*base.Response, *gortsplib.ServerStream, error) {
	me.logger.Debugf("SETUP request %+v", ctx)
	clt := getClient(ctx.Conn)
	if clt.stream == nil {
		return &base.Response{
			StatusCode:    base.StatusBadRequest,
			StatusMessage: "you must call DESCRIBE before SETUP",
		}, nil, nil
	}

	trackIDStr := strings.TrimPrefix(strings.TrimPrefix(ctx.Request.URL.Path, clt.path), "/trackID=")
	trackID64, err := strconv.ParseUint(trackIDStr, 10, 32)
	if err != nil {
		return &base.Response{
			StatusCode:    base.StatusBadRequest,
			StatusMessage: "invalid track ID",
			Body:          []byte(fmt.Sprintf("invalid track ID: %s", trackIDStr)),
		}, nil, nil
	}
	trackID := int(trackID64)
	if trackID > len(clt.stream.Desc.Medias) {
		return &base.Response{
			StatusCode:    base.StatusBadRequest,
			StatusMessage: "track ID out of range",
			Body:          []byte(fmt.Sprintf("track ID out of range: %d", trackID)),
		}, nil, nil
	}
	media := clt.stream.Desc.Medias[trackID]

	url := clt.config.URL
	me.logger.Debugf("Proxying SETUP request to %s", url.String())
	res, err := clt.Setup((*base.URL)(&url), media, 0, 0)
	if err != nil {
		me.logger.WithError(err).Errorf("error proxying SETUP request")
		return wrapError(errGatewayTimeout, err), nil, nil
	}
	return res, clt.stream, nil
}

func (me *server) OnPlay(ctx *gortsplib.ServerHandlerOnPlayCtx) (*base.Response, error) {
	me.logger.Debugf("PLAY request %+v", ctx)
	clt := getClient(ctx.Conn)
	if clt.stream == nil {
		return &base.Response{
			StatusCode:    base.StatusBadRequest,
			StatusMessage: "you must call DESCRIBE before PLAY",
		}, nil
	}

	clt.OnPacketRTCPAny(func(m *description.Media, p rtcp.Packet) {
		err := clt.stream.WritePacketRTCP(m, p)
		if err != nil {
			me.logger.WithError(err).Errorf("error writing RTCP packet from client")
		}
	})
	clt.OnPacketRTPAny(func(m *description.Media, f format.Format, p *rtp.Packet) {
		err := clt.stream.WritePacketRTP(m, p)
		if err != nil {
			me.logger.WithError(err).Errorf("error writing RTP packet from client")
		}
	})

	rngStr, found := ctx.Request.Header["Range"]
	if !found {
		return &base.Response{
			StatusCode:    base.StatusBadRequest,
			StatusMessage: "Range header missing",
			Body:          []byte("Range header missing"),
		}, nil
	}
	var rng headers.Range
	err := rng.Unmarshal(rngStr)
	if err != nil {
		return &base.Response{
			StatusCode:    base.StatusBadRequest,
			StatusMessage: "Range header invalid",
			Body:          []byte(err.Error()),
		}, nil
	}

	url := clt.config.URL
	me.logger.Debugf("Proxying PLAY request to %s", url.String())
	res, err := clt.Play(&rng)
	if err != nil {
		me.logger.WithError(err).Errorf("error proxying PLAY request")
		return wrapError(errGatewayTimeout, err), nil
	}
	return res, nil
}

func (me *server) OnSetParameter(ctx *gortsplib.ServerHandlerOnSetParameterCtx) (*base.Response, error) {
	me.logger.Debugf("SET_PARAMETER request %+v", ctx)
	return errNotImplemented, nil
}

func (me *server) OnGetParameter(ctx *gortsplib.ServerHandlerOnGetParameterCtx) (*base.Response, error) {
	me.logger.Debugf("GET_PARAMETER request %+v", ctx)
	return errNotImplemented, nil
}

func (me *server) OnPause(ctx *gortsplib.ServerHandlerOnPauseCtx) (*base.Response, error) {
	me.logger.Debugf("PAUSE request %+v", ctx)
	clt := getClient(ctx.Conn)

	url := clt.config.URL
	me.logger.Debugf("Proxying PAUSE request to %s", url.String())
	res, err := clt.Pause()
	if err != nil {
		me.logger.WithError(err).Errorf("error proxying PAUSE request")
		return wrapError(errGatewayTimeout, err), nil
	}
	return res, nil
}

func (me *server) OnRecord(ctx *gortsplib.ServerHandlerOnRecordCtx) (*base.Response, error) {
	me.logger.Debugf("RECORD request %+v", ctx)
	clt := getClient(ctx.Conn)

	url := clt.config.URL
	me.logger.Debugf("Proxying RECORD request to %s", url.String())
	res, err := clt.Record()
	if err != nil {
		me.logger.WithError(err).Errorf("error proxying RECORD request")
		return wrapError(errGatewayTimeout, err), nil
	}
	return res, nil
}

func (me *server) OnAnnounce(ctx *gortsplib.ServerHandlerOnAnnounceCtx) (*base.Response, error) {
	me.logger.Debugf("ANNOUNCE request %+v", ctx)
	clt := getClient(ctx.Conn)
	if clt.stream == nil {
		return &base.Response{
			StatusCode:    base.StatusBadRequest,
			StatusMessage: "you must call DESCRIBE before ANNOUNCE",
		}, nil
	}

	url := clt.config.URL
	me.logger.Debugf("Proxying ANNOUNCE request to %s", url.String())
	res, err := clt.Announce((*base.URL)(&url), clt.stream.Desc)
	if err != nil {
		me.logger.WithError(err).Errorf("error proxying ANNOUNCE request")
		return wrapError(errGatewayTimeout, err), nil
	}
	return res, nil
}
