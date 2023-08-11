package response

import "errors"

var (
	ErrResponseVersionNotSupported = errors.New("gateway[response]: version not supported")
)
