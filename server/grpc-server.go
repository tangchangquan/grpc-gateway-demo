package server

import (
	"log"
	"net"

	pb "github.com/helloworlde/grpc-gateway/proto/api"
	"github.com/helloworlde/grpc-gateway/service"
	"google.golang.org/grpc"
)

var helloService = service.HelloService{}

func StartGrpcServer() {
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalln("Listen gRPC port failed: ", err)
	}

	server := grpc.NewServer()
	pb.RegisterHelloServiceServer(server, &helloService)

	log.Println("Start gRPC Server on 0.0.0.0:9090")
	err = server.Serve(listener)
	if err != nil {
		log.Fatalln("Start gRPC Server failed: ", err)
	}

}
