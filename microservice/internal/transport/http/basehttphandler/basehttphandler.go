package basehttphandler

import (
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/pkg"
	"log/slog"
	"time"
)

// BaseHttpHandler is the base handler for http requests.
type BaseHttpHandler struct {
	*pkg.Packages
	Logger        *slog.Logger
	CancelTimeout time.Duration
}

// Option is the option type for http handler.
type Option func(*BaseHttpHandler)

// WithContextTimeout sets the cancel timeout option.
func WithContextTimeout(timeout time.Duration) Option {
	return func(h *BaseHttpHandler) {
		h.CancelTimeout = timeout
	}
}

// WithLogger sets the logger option.
func WithLogger(logger *slog.Logger) Option {
	return func(h *BaseHttpHandler) {
		h.Logger = logger
	}
}

// WithPackages sets the packages option.
func WithPackages(packages *pkg.Packages) Option {
	return func(h *BaseHttpHandler) {
		h.Packages = packages
	}
}

// New creates a new http handler with the given options.
func New(opts ...Option) *BaseHttpHandler {
	h := &BaseHttpHandler{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}
