# Lamway

Is a manage layer implementation for AWS Lambda functions using API Gateway v1 or v2.

It provides a way to use go base http.Handler instances to manage all needed HTTP paths on the same lambda.

## Requirements

- Go >= 1.20

## Installation 

```bash
go get github.com/danteay/lamway
```

## Example

- [API Gateway v1](https://github.com/danteay/lamway/tree/main/examples/api-gateway-v1)
- [API Gateway v2](https://github.com/danteay/lamway/tree/main/examples/api-gateway-v2)
- [Gin](https://github.com/danteay/lamway/tree/main/examples/gin)

## Credits

- Based on [apex/gateway](https://github.com/apex/gateway)
