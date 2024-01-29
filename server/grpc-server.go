package server

import (
	"log"
	"net"

	pb "github.com/helloworlde/grpc-gateway/proto/api"
	"github.com/helloworlde/grpc-gateway/service"
	"google.golang.org/grpc"
)

var helloService = service.HelloService{}
var grpcAddress = "127.0.0.1"
var grpcport = 9090

func StartGrpcServer() {
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalln("Listen gRPC port failed: ", err)
	}

	//注册到consul
	if err := consulGrpc(); err != nil {
		log.Fatalln("consulGrpc failed: ", err)
	}

	server := grpc.NewServer()
	pb.RegisterHelloServiceServer(server, &helloService)

	log.Println("Start gRPC Server on 0.0.0.0:9090")
	err = server.Serve(listener)
	if err != nil {
		log.Fatalln("Start gRPC Server failed: ", err)
	}

}
