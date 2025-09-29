package response

type APIGatewayResponse struct {
	StatusCode        int                 `json:"statusCode"`
	Headers           map[string]string   `json:"headers"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	Body              string              `json:"body"`
	IsBase64Encoded   bool                `json:"isBase64Encoded,omitempty"`

	// Just for APIGateway v2
	Cookies []string `json:"cookies"`
}

func (agr APIGatewayResponse) ToV1Map() map[string]any {
	return map[string]any{
		"statusCode":        agr.StatusCode,
		"headers":           agr.Headers,
		"multiValueHeaders": agr.MultiValueHeaders,
		"body":              agr.Body,
		"isBase64Encoded":   agr.IsBase64Encoded,
	}
}

func (agr APIGatewayResponse) ToV2Map() map[string]any {
	return map[string]any{
		"statusCode":        agr.StatusCode,
		"headers":           agr.Headers,
		"multiValueHeaders": agr.MultiValueHeaders,
		"body":              agr.Body,
		"isBase64Encoded":   agr.IsBase64Encoded,
		"cookies":           agr.Cookies,
	}
}
