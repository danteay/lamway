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
		assert.JSONEq(t, `{"body":"Hello World from Go\n", "cookies":null, "headers":{"Content-Type":"text/plain; charset=utf8"}, "multiValueHeaders":{}, "statusCode":200}`, string(res))
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
		assert.JSONEq(t, `{"body":"Hello World from Go\n", "cookies":null, "headers":{"Content-Type":"text/plain; charset=utf8"}, "multiValueHeaders":{}, "statusCode":200}`, string(res))
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
		assert.Equal(t, gw.defaultResponse, res)
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
		assert.Equal(t, gw.defaultResponse, res)
	})
}
