package proxy

import (
	"context"
	"sync"
	"time"

	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
)

type pubSub struct {
	mu   sync.Mutex
	subs map[string][]chan error
}

func newPubSub() *pubSub {
	return &pubSub{subs: make(map[string][]chan error)}
}

func (me *pubSub) subscribe(path string) chan error {
	ch := make(chan error, 1)
	me.mu.Lock()
	defer me.mu.Unlock()
	me.subs[path] = append(me.subs[path], ch)
	return ch
}

func (me *pubSub) unsubscribe(path string, ch chan error) {
	me.mu.Lock()
	defer me.mu.Unlock()
	subs := me.subs[path]
	for i, s := range subs {
		if s == ch {
			me.subs[path] = append(subs[:i], subs[i+1:]...)
			break
		}
	}
}

func (me *pubSub) publish(path string, err error) {
	me.mu.Lock()
	defer me.mu.Unlock()
	for _, ch := range me.subs[path] {
		select {
		case ch <- err:
		default:
		}
	}
}

func (me *server) runMonitor(ctx context.Context) {
	me.probeAll()
	ticker := time.NewTicker(upstreamMonitorInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			me.probeAll()
		}
	}
}

func (me *server) probeAll() {
	for path, stream := range me.streams {
		go me.probeUpstream(path, stream)
	}
}

func (me *server) probeUpstream(path string, cfg *StreamConfig) {
	url := cfg.URL
	probe := gortsplib.Client{
		Scheme:       url.Scheme,
		Host:         url.Host,
		UserAgent:    userAgent,
		ReadTimeout:  connTimeout,
		WriteTimeout: connTimeout,
	}
	err := probe.Start2()
	if err != nil {
		me.updateConnErr(path, err)
		return
	}
	defer probe.Close()
	_, _, err = probe.Describe((*base.URL)(&url))
	me.updateConnErr(path, err)
}

func (me *server) updateConnErr(path string, newErr error) {
	me.thumbsMutex.Lock()
	thumb, ok := me.thumbs[path]
	if !ok {
		me.thumbsMutex.Unlock()
		return
	}
	oldErr := thumb.connErr
	thumb.connErr = newErr
	me.thumbsMutex.Unlock()

	wasOffline := oldErr != nil
	isOfflineNow := newErr != nil
	if wasOffline == isOfflineNow {
		return
	}
	if newErr == nil {
		me.logger.Infof("stream %q came back online", path)
	} else {
		me.logger.WithError(newErr).Warnf("stream %q went offline", path)
	}
	me.pubSub.publish(path, newErr)
}
