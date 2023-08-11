package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"

	"github.com/danteay/lamway"
)

var r *gin.Engine

func init() {
	r = gin.Default()

	r.GET("/", helloWorld)
}

func main() {
	gw := lamway.New[events.APIGatewayProxyRequest](
		lamway.WithHTTPHandler(r),
	)

	log.Fatal(gw.Start())
}

func helloWorld(gc *gin.Context) {
	gc.JSON(http.StatusOK, gin.H{
		"message": "Hello World from Go - method " + gc.Request.Method,
	})
}
