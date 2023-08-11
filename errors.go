package lamway

import "errors"

var (
	ErrInvalidAPIGatewayRequest = errors.New("gateway: invalid APIGateway request struct configured")
)
