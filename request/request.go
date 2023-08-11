package request

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func NewV1(ctx context.Context, evt events.APIGatewayProxyRequest) (*http.Request, error) {
	ri, err := newAPIGatewayV1RequestInfo(evt)
	if err != nil {
		return nil, err
	}

	return ri.toRequest(ctx)
}

func NewV2(ctx context.Context, evt events.APIGatewayV2HTTPRequest) (*http.Request, error) {
	ri := newAPIGatewayV2RequestInfo(evt)
	return ri.toRequest(ctx)
}
