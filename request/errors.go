package request

import "errors"

var (
	ErrParsingPathFailed   = errors.New("gateway[request]: parsing path failed")
	ErrDecodingBase64Body  = errors.New("gateway[request]: decoding base64 body")
	ErrFailToCreateRequest = errors.New("gateway[request]: fail to create request")
)
