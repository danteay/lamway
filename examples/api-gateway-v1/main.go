package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/danteay/lamway"
)

func init() {
	http.HandleFunc("/", helloWorld)
}

func main() {
	gw := lamway.New[events.APIGatewayProxyRequest]()
	log.Fatal(gw.Start())
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World from Go - method " + r.Method))
}
