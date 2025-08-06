package proxy

import (
	"context"
	"fmt"

	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
	"github.com/sirupsen/logrus"
)

var (
	_ gortsplib.ServerHandlerOnConnOpen         = &server{}
	_ gortsplib.ServerHandlerOnConnClose        = &server{}
	_ gortsplib.ServerHandlerOnSessionOpen      = &server{}
	_ gortsplib.ServerHandlerOnSessionClose     = &server{}
	_ gortsplib.ServerHandlerOnRequest          = &server{}
	_ gortsplib.ServerHandlerOnResponse         = &server{}
	_ gortsplib.ServerHandlerOnDescribe         = &server{}
	_ gortsplib.ServerHandlerOnAnnounce         = &server{}
	_ gortsplib.ServerHandlerOnSetup            = &server{}
	_ gortsplib.ServerHandlerOnPlay             = &server{}
	_ gortsplib.ServerHandlerOnRecord           = &server{}
	_ gortsplib.ServerHandlerOnPause            = &server{}
	_ gortsplib.ServerHandlerOnGetParameter     = &server{}
	_ gortsplib.ServerHandlerOnSetParameter     = &server{}
	_ gortsplib.ServerHandlerOnPacketsLost      = &server{}
	_ gortsplib.ServerHandlerOnDecodeError      = &server{}
	_ gortsplib.ServerHandlerOnStreamWriteError = &server{}
)

const userAgent = "RTSP-Fixer"

type server struct {
	rtsp    *gortsplib.Server
	logger  *logrus.Logger
	streams map[string]*StreamConfig
}

func NewServer(logger *logrus.Logger, port string, streams []StreamConfig) *server {
	streamsMap := make(map[string]*StreamConfig, len(streams))
	for _, stream := range streams {
		streamsMap[fmt.Sprintf("/%s", stream.Name)] = &stream
	}
	me := &server{
		logger:  logger,
		streams: streamsMap,
	}
	me.rtsp = &gortsplib.Server{
		Handler:     me,
		RTSPAddress: fmt.Sprintf(":%s", port),
	}
	return me
}

func (me *server) Run(ctx context.Context) error {
	go func() {
		select {
		case <-ctx.Done():
			me.logger.Infof("stopping server...")
			me.rtsp.Close()
		}
	}()
	return me.rtsp.StartAndWait()
}

func (me *server) OnConnOpen(ctx *gortsplib.ServerHandlerOnConnOpenCtx) {
	me.logger.Debugf("conn opened")
}

func (me *server) OnRequest(ctx *gortsplib.ServerConn, req *base.Request) {
	me.logger.Debugf("request %+v", req)
	if ctx.UserData() != nil {
		return
	}
	client, err := me.createClientFor(req.URL.Path)
	if err != nil {
		me.logger.WithError(err).Errorf("failed to create client for %q", req.URL.Path)
		ctx.Close()
		return
	}
	me.logger.Debugf("created client for %q", req.URL.Path)
	ctx.SetUserData(client)
}

func (me *server) OnSessionOpen(ctx *gortsplib.ServerHandlerOnSessionOpenCtx) {
	me.logger.Debugf("session opened")
}

func (me *server) OnSessionClose(ctx *gortsplib.ServerHandlerOnSessionCloseCtx) {
	me.logger.Debugf("session closed")
}

func (me *server) OnResponse(ctx *gortsplib.ServerConn, resp *base.Response) {
	me.logger.Debugf("response %+v", resp)
	resp.Header["Server"] = base.HeaderValue{userAgent}
}

func (me *server) OnConnClose(ctx *gortsplib.ServerHandlerOnConnCloseCtx) {
	me.logger.WithError(ctx.Error).Info("conn closed")
	client, err := getOptionalClient(ctx.Conn)
	if err != nil {
		return
	}
	if client.stream != nil {
		client.stream.Close()
	}
	client.Close()
}

func (me *server) OnStreamWriteError(ctx *gortsplib.ServerHandlerOnStreamWriteErrorCtx) {
	me.logger.WithError(ctx.Error).Error("stream write error")
}

func (me *server) OnDecodeError(ctx *gortsplib.ServerHandlerOnDecodeErrorCtx) {
	me.logger.WithError(ctx.Error).Error("decode error")
}

func (me *server) OnPacketsLost(ctx *gortsplib.ServerHandlerOnPacketsLostCtx) {
	me.logger.Errorf("lost %d packets", ctx.Lost)
}
