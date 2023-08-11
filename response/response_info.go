package response

import "github.com/aws/aws-lambda-go/events"

type responseInfo struct {
	statusCode        int
	headers           map[string]string
	multiValueHeaders map[string][]string
	isBase64Encoded   bool
	body              string
	cookies           []string
}

func (ri responseInfo) toAPIGatewayV1Response() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode:        ri.statusCode,
		Headers:           ri.headers,
		MultiValueHeaders: ri.multiValueHeaders,
		Body:              ri.body,
		IsBase64Encoded:   ri.isBase64Encoded,
	}
}

func (ri responseInfo) toAPIGatewayV2Response() events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode:        ri.statusCode,
		Headers:           ri.headers,
		MultiValueHeaders: ri.multiValueHeaders,
		Body:              ri.body,
		IsBase64Encoded:   ri.isBase64Encoded,
		Cookies:           ri.cookies,
	}
}
