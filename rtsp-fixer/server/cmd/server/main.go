package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/LouisBrunner/ha-addons/rtsp-fixer/server/pkg/proxy"
	"github.com/sirupsen/logrus"
)

type config struct {
	port    string
	streams []proxy.StreamConfig
}

func parseStreams(streamsStr string) ([]proxy.StreamConfig, error) {
	var streams []proxy.StreamConfig
	dec := json.NewDecoder(strings.NewReader(streamsStr))
	for {
		var stream proxy.StreamConfig
		err := dec.Decode(&stream)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		streams = append(streams, stream)
	}
	return streams, nil
}

func parseConfig(logger *logrus.Logger) (*config, error) {
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	}
	if os.Getenv("DEBUG") == "true" {
		logger.SetLevel(logrus.DebugLevel)
		logger.SetReportCaller(true)
	}
	logger.SetFormatter(formatter)

	logger.Infof("starting server")
	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("internal error: PORT environment variable is not set")
	}
	streamsStr := os.Getenv("STREAMS")
	if streamsStr == "" {
		return nil, fmt.Errorf("no streams have been provided, edit your configuration")
	}
	streams, err := parseStreams(streamsStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing streams %q: %w", streamsStr, err)
	}
	if len(streams) == 0 {
		return nil, fmt.Errorf("no streams have been provided, edit your configuration")
	}

	for _, stream := range streams {
		logger.Infof("proxying %s to %s (fix TCP in transport: %v)", stream.Name, stream.URL.String(), stream.FixForceTCPInTransport)
	}

	return &config{
		port:    port,
		streams: streams,
	}, nil
}

func work(ctx context.Context, logger *logrus.Logger) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctx, cleanup := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cleanup()

	cfg, err := parseConfig(logger)
	if err != nil {
		return err
	}

	logger.Infof("starting server on port %s", cfg.port)
	server := proxy.NewServer(logger, cfg.port, cfg.streams)
	return server.Run(ctx)
}

func main() {
	logger := logrus.New()
	err := work(context.Background(), logger)
	if err != nil {
		logger.WithError(err).Error("fatal error, exiting")
		os.Exit(1)
	}
}
