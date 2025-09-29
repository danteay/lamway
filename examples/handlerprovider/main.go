package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"

	"github.com/danteay/lamway"
)

func main() {
	gw := lamway.New[events.APIGatewayProxyRequest](
		// this wraps the http.Handler initialization with the context provided by the lambda.Start() call
		// to be able to propagate it forward to the handler
		lamway.WithHandlerProvider(func(_ context.Context) http.Handler {
			r := gin.Default()

			r.GET("/", helloWorld)

			return r
		}),
	)

	log.Fatal(gw.Start())
}

func helloWorld(gc *gin.Context) {
	gc.JSON(http.StatusOK, gin.H{
		"message": "Hello World from Go - method " + gc.Request.Method,
	})
}
