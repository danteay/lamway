package lamway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"

	"github.com/danteay/lamway/request"
)

const testPath = "/pets/luna"

func hello(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Custom-Header", "custom-value")

	_, _ = fmt.Fprintln(w, "Hello World from Go")
}

func TestGateway_invoke(t *testing.T) {
	t.Run("should execute api gateway v1", func(t *testing.T) {
		evt := events.APIGatewayProxyRequest{
			Path:       testPath,
			HTTPMethod: http.MethodPost,
		}

		gw := New[events.APIGatewayProxyRequest](WithHTTPHandler(http.HandlerFunc(hello)))

		payload, err := gw.invoke(context.Background(), evt)

		res, err := json.Marshal(payload)
		if err != nil {
			assert.Fail(t, "can't marshal payload", err)
		}

		assert.NoError(t, err)
		assert.JSONEq(t, `{"body":"Hello World from Go\n", "headers":{"Content-Type":"text/plain; charset=utf8", "Custom-Header":"custom-value"}, "isBase64Encoded":false, "multiValueHeaders":{}, "statusCode":200}`, string(res))
	})

	t.Run("should execute api gateway v2", func(t *testing.T) {
		evt := events.APIGatewayV2HTTPRequest{
			RawPath: testPath,
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Method: http.MethodPost,
					Path:   testPath,
				},
			},
		}

		gw := New[events.APIGatewayV2HTTPRequest](WithHTTPHandler(http.HandlerFunc(hello)))

		payload, err := gw.invoke(context.Background(), evt)

		res, err := json.Marshal(payload)
		if err != nil {
			assert.Fail(t, "can't marshal payload", err)
		}

		assert.NoError(t, err)
		assert.JSONEq(t, `{"body":"Hello World from Go\n", "cookies":null, "headers":{"Content-Type":"text/plain; charset=utf8", "Custom-Header":"custom-value"}, "isBase64Encoded":false, "multiValueHeaders":{}, "statusCode":200}`, string(res))
	})

	t.Run("should get error bay error parsing path on v1", func(t *testing.T) {
		evt := events.APIGatewayProxyRequest{
			Path:       testPath + string(rune(0x7f)),
			HTTPMethod: http.MethodPost,
		}

		gw := New[events.APIGatewayProxyRequest](WithHTTPHandler(http.HandlerFunc(hello)))

		res, err := gw.invoke(context.Background(), evt)

		assert.Error(t, err)
		assert.ErrorIs(t, err, request.ErrParsingPathFailed)
		assert.Equal(t, gw.defaultResponse.ToV1Map(), res)
	})

	t.Run("should get error bay error parsing path on v2", func(t *testing.T) {
		evt := events.APIGatewayV2HTTPRequest{
			RawPath: testPath + string(rune(0x7f)),
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Method: http.MethodPost,
					Path:   testPath + string(rune(0x7f)),
				},
			},
		}

		gw := New[events.APIGatewayV2HTTPRequest](WithHTTPHandler(http.HandlerFunc(hello)))

		res, err := gw.invoke(context.Background(), evt)

		assert.Error(t, err)
		assert.ErrorIs(t, err, request.ErrParsingPathFailed)
		assert.Equal(t, gw.defaultResponse.ToV2Map(), res)
	})
}

func TestGateway_WithHandlerProvider(t *testing.T) {
	// Ensure that the handler provider is used to lazily construct the handler
	// and that it is only called once even across multiple invocations.

	t.Run("v1 provider is used and called once", func(t *testing.T) {
		var called int
		provider := func(_ context.Context) http.Handler {
			called++
			return http.HandlerFunc(hello)
		}

		gw := New[events.APIGatewayProxyRequest](WithHandlerProvider(provider))

		evt := events.APIGatewayProxyRequest{
			Path:       testPath,
			HTTPMethod: http.MethodPost,
		}

		// The first invocation should initialize the handler via provider
		payload1, err1 := gw.invoke(context.Background(), evt)
		res1, merr1 := json.Marshal(payload1)
		if merr1 != nil {
			assert.Fail(t, "can't marshal payload", merr1)
		}
		assert.NoError(t, err1)
		assert.JSONEq(t, `{"body":"Hello World from Go\n", "headers":{"Content-Type":"text/plain; charset=utf8", "Custom-Header":"custom-value"}, "isBase64Encoded":false, "multiValueHeaders":{}, "statusCode":200}`, string(res1))

		// The second invocation should reuse the same handler (provider not called again)
		payload2, err2 := gw.invoke(context.Background(), evt)
		res2, merr2 := json.Marshal(payload2)
		if merr2 != nil {
			assert.Fail(t, "can't marshal payload", merr2)
		}
		assert.NoError(t, err2)
		assert.JSONEq(t, `{"body":"Hello World from Go\n", "headers":{"Content-Type":"text/plain; charset=utf8", "Custom-Header":"custom-value"}, "isBase64Encoded":false, "multiValueHeaders":{}, "statusCode":200}`, string(res2))

		assert.Equal(t, 1, called, "handler provider should be called exactly once")
	})

	t.Run("v2 provider is used and called once", func(t *testing.T) {
		var called int
		provider := func(_ context.Context) http.Handler {
			called++
			return http.HandlerFunc(hello)
		}

		gw := New[events.APIGatewayV2HTTPRequest](WithHandlerProvider(provider))

		evt := events.APIGatewayV2HTTPRequest{
			RawPath: testPath,
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Method: http.MethodPost,
					Path:   testPath,
				},
			},
		}

		payload1, err1 := gw.invoke(context.Background(), evt)
		res1, merr1 := json.Marshal(payload1)
		if merr1 != nil {
			assert.Fail(t, "can't marshal payload", merr1)
		}
		assert.NoError(t, err1)
		assert.JSONEq(t, `{"body":"Hello World from Go\n", "cookies":null, "headers":{"Content-Type":"text/plain; charset=utf8", "Custom-Header":"custom-value"}, "isBase64Encoded":false, "multiValueHeaders":{}, "statusCode":200}`, string(res1))

		payload2, err2 := gw.invoke(context.Background(), evt)
		res2, merr2 := json.Marshal(payload2)
		if merr2 != nil {
			assert.Fail(t, "can't marshal payload", merr2)
		}
		assert.NoError(t, err2)
		assert.JSONEq(t, `{"body":"Hello World from Go\n", "cookies":null, "headers":{"Content-Type":"text/plain; charset=utf8", "Custom-Header":"custom-value"}, "isBase64Encoded":false, "multiValueHeaders":{}, "statusCode":200}`, string(res2))

		assert.Equal(t, 1, called, "handler provider should be called exactly once")
	})
}
