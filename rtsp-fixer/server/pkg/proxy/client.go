package proxy

import (
	"fmt"
	"strings"

	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
)

type client struct {
	srv *server
	gortsplib.Client
	path   string
	config *StreamConfig
	stream *gortsplib.ServerStream
}

func (me *server) createClientFor(path string) (*client, error) {
	stream, ok := me.streams[path]
	if !ok {
		return nil, fmt.Errorf("stream %q not found", path)
	}

	clt := &client{
		srv: me,
		Client: gortsplib.Client{
			Scheme:    stream.URL.Scheme,
			Host:      stream.URL.Host,
			UserAgent: userAgent,
		},
		path:   path,
		config: stream,
	}
	clt.Client.OnRequest = clt.onRequest
	clt.Client.OnResponse = clt.onResponse
	clt.Client.OnTransportSwitch = clt.onTransportSwitch
	clt.Client.OnPacketsLost = clt.onPacketsLost
	clt.Client.OnDecodeError = clt.onDecodeError
	err := clt.Start2()
	if err != nil {
		return nil, err
	}

	return clt, nil
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
