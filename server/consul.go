package server

import (
	"log"

	consul "github.com/helloworlde/grpc-gateway/consul"
)

func init() {
	consul.Init(nil)
}

func consulGrpc() error {

	// 创建gRPC服务器实例
	if err := consul.Register("hello", grpcAddress, grpcport, "hellotag"); err != nil {
		log.Fatalf("consul err:%v", err.Error())
		return err
	}
	return nil

}
