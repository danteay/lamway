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

## Development Tasks

This project uses go-task to run common tasks. After installing go-task (or using the provided Nix flake), you can run:

- task install      # Download Go module dependencies
- task pre-commit   # Install pre-commit hooks
- task vet          # Run go vet checks
- task lint         # Run revive linter (depends on vet)
- task format       # Format code using gofmt
- task test         # Run unit tests with race and coverage
- task cover        # Open HTML coverage report (runs tests first)

If you use Nix, enter the dev shell with:

```bash
nix develop
```

Otherwise, install go-task from https://taskfile.dev.

## Credits

- Based on [apex/gateway](https://github.com/apex/gateway)
