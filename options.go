package lamway

import (
	"net/http"
)

type options struct {
	httpHandler     http.Handler
	handlerProvider HandlerProvider
	decorators      []Decorator
	defaultHeaders  map[string]string
	defaultErrorRes string
	logger          Logger
}

// Option is a functional option for configuring the gateway.
type Option func(*options)

// WithHTTPHandler sets the provided http.Handler to the options' configuration.
func WithHTTPHandler(h http.Handler) Option {
	return func(o *options) {
		o.httpHandler = h
	}
}

// WithLogger sets the provided Logger instance to the options' configuration.
func WithLogger(l Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}

// WithDecorator adds a Decorator function to the options for processing handlers.
func WithDecorator(d Decorator) Option {
	return func(o *options) {
		if o.decorators == nil {
			o.decorators = make([]Decorator, 0)
		}

		o.decorators = append(o.decorators, d)
	}
}

// WithDefaultErrorHeaders sets a default error response header to be used in the options if the provided map is not empty.
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

// WithDefaultErrorResponse sets a default error response string to be used in the options if the provided string is not empty.
func WithDefaultErrorResponse(res string) Option {
	return func(o *options) {
		if res == "" {
			return
		}

		o.defaultErrorRes = res
	}
}

// WithHandlerProvider sets a custom HandlerProvider to lazily initialize the HTTP handler with context propagation.
func WithHandlerProvider(p HandlerProvider) Option {
	return func(o *options) {
		o.handlerProvider = p
	}
}
