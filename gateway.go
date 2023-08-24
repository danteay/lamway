package lamway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/danteay/lamway/request"
	"github.com/danteay/lamway/response"
	"github.com/danteay/lamway/types"
)

type Logger interface {
	Debugf(format string, args ...any)
}

// Gateway wrap a http handler to enable use as a lambda.Handler
type Gateway[T any] struct {
	handler         http.Handler
	decorators      []types.Decorator
	defaultResponse types.APIGatewayResponse
	logger          Logger
}

// New creates a gateway using the provided http.Handler enabling use in existing aws-lambda-go
// projects
func New[T any](opts ...Option) *Gateway[T] {
	gatewayOpts := options{
		httpHandler:     http.DefaultServeMux,
		defaultHeaders:  map[string]string{"Content-Type": "application/json"},
		defaultErrorRes: `{"message": "Error processing request"}`,
	}

	for _, opt := range opts {
		opt(&gatewayOpts)
	}

	return &Gateway[T]{
		handler:    gatewayOpts.httpHandler,
		decorators: gatewayOpts.decorators,
		logger:     gatewayOpts.logger,
		defaultResponse: types.APIGatewayResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    gatewayOpts.defaultHeaders,
			Body:       gatewayOpts.defaultErrorRes,
		},
	}
}

// GetInvoker returns the function that will be invoked by the lambda.Start call in main function. This funtion will be
// decorated or not depending on the options passed to the New function.
func (gw *Gateway[T]) GetInvoker() any {
	var worker any = gw.invoke

	if len(gw.decorators) > 0 {
		for _, decorator := range gw.decorators {
			worker = decorator(worker)
		}
	}

	return worker
}

func (gw *Gateway[T]) Start() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
				return
			}

			err = fmt.Errorf("%v", r)
		}
	}()

	lambda.Start(gw.GetInvoker())

	return nil
}

func (gw *Gateway[T]) invoke(ctx context.Context, evt T) (map[string]any, error) {
	aux := any(evt)

	switch v := aux.(type) {
	case events.APIGatewayProxyRequest:
		gw.logDebug("[id:%s] v1 request: %+v", v.RequestContext.RequestID, aux)

		res, err := gw.handlerV1(ctx, v)
		apiRes := res.ToV1Map()

		gw.logDebug("[id:%s] v1 response: %+v", v.RequestContext.RequestID, apiRes)

		return apiRes, err
	case events.APIGatewayV2HTTPRequest:
		gw.logDebug("[id:%s] v2 request: %+v", v.RequestContext.RequestID, aux)

		res, err := gw.handlerV2(ctx, v)
		apiRes := res.ToV2Map()

		gw.logDebug("[id:%s] v2 response: %+v", v.RequestContext.RequestID, apiRes)

		return apiRes, err
	default:
		return gw.defaultResponse.ToV1Map(), ErrInvalidAPIGatewayRequest
	}
}

func (gw *Gateway[T]) handlerV1(ctx context.Context, evt events.APIGatewayProxyRequest) (types.APIGatewayResponse, error) {
	r, err := request.NewV1(ctx, evt)
	if err != nil {
		return gw.defaultResponse, err
	}

	w := response.New()

	gw.handler.ServeHTTP(w, r)

	return w.End(), nil
}

func (gw *Gateway[T]) handlerV2(ctx context.Context, evt events.APIGatewayV2HTTPRequest) (types.APIGatewayResponse, error) {
	r, err := request.NewV2(ctx, evt)
	if err != nil {
		return gw.defaultResponse, err
	}

	w := response.New()

	gw.handler.ServeHTTP(w, r)

	return w.End(), nil
}

func (gw *Gateway[T]) logDebug(format string, args ...any) {
	if gw.logger != nil {
		gw.logger.Debugf(format, args...)
	}
}
