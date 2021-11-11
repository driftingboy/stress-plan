package http

import (
	"net/http"
	"stress-plan/logger"
	"time"

	"golang.org/x/net/http2"
)

// for new sender
type Option func(*Sender)

func WithTTL(ttl time.Duration) Option {
	return func(s *Sender) {
		s.ttl = ttl
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(s *Sender) {
		s.logger = logger
	}
}

func WithHttp2() Option {
	return func(s *Sender) {
		if s.cli != nil && s.cli.Transport != nil {
			if tr, ok := s.cli.Transport.(*http.Transport); ok {
				http2.ConfigureTransport(tr)
				s.cli.Transport = tr
			}
		}
	}
}

func WithDisableKeepAlive() Option {
	return func(s *Sender) {
		if s.cli != nil && s.cli.Transport != nil {
			if tr, ok := s.cli.Transport.(*http.Transport); ok {
				tr.DisableKeepAlives = true
				s.cli.Transport = tr
			}
		}
	}
}

type TrConfig struct {
}

func WithTransPort(tc *TrConfig) Option {
	return func(s *Sender) {}
}
