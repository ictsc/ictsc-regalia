package pgxutil

import "github.com/XSAM/otelsql"

type OpenOption interface {
	applyOpenOption(cfg *openConfig)
}

type openOptionFunc func(*openConfig)

func (f openOptionFunc) applyOpenOption(cfg *openConfig) {
	f(cfg)
}

func WithOTel(enabled bool) OpenOption {
	return openOptionFunc(func(cfg *openConfig) {
		cfg.enableOTel = enabled
	})
}

func WithOTelOptions(opts ...otelsql.Option) OpenOption {
	return openOptionFunc(func(cfg *openConfig) {
		cfg.otelOptions = append(cfg.otelOptions, opts...)
	})
}
