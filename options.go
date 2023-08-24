package lamway

import (
	"net/http"

	"github.com/danteay/lamway/types"
)

type options struct {
	httpHandler     http.Handler
	decorators      []types.Decorator
	defaultHeaders  map[string]string
	defaultErrorRes string
	logger          Logger
}

type Option func(*options)

func WithHTTPHandler(h http.Handler) Option {
	return func(o *options) {
		o.httpHandler = h
	}
}

func WithLogger(l Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}

func WithDecorator(d types.Decorator) Option {
	return func(o *options) {
		if o.decorators == nil {
			o.decorators = make([]types.Decorator, 0)
		}

		o.decorators = append(o.decorators, d)
	}
}

func WithDefaultErrorHeaders(headers map[string]string) Option {
	return func(o *options) {
		if headers == nil {
			return
		}

		for k, v := range headers {
			o.defaultHeaders[k] = v
		}
	}
}

func WithDefaultErrorResponse(res string) Option {
	return func(o *options) {
		if res == "" {
			return
		}

		o.defaultErrorRes = res
	}
}
