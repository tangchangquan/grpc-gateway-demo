package main

import (
	"github.com/helloworlde/grpc-gateway/server"
)

func main() {
	go server.StartGrpcServer()
	server.StartGwServer()
}
