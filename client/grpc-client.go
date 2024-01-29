package main

import (
	"context"
	"flag"
	helloworld "github.com/helloworlde/grpc-gateway/proto/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"log"
)

const (
	defaultName = "tcq"
)

//var (
//	addr = flag.String("addr", "127.0.0.1:9090", "the address to connect to")
//)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		grpclog.Fatal(err)
	}

}

func run() error {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := grpc.Dial(
		// consul://192.168.193.128:8500 consul地址
		// test-serve 拉取的服务名
		// wait=14s 等待时间
		// tag=manual 筛选条件
		// 底层就是利用grpc-consul-resolver将参数解析成HTTP请求获取对应的服务
		"consul://127.0.0.1:8500/hello?wait=5s&tag=hellotag",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"hello": "helloService"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := helloworld.NewHelloServiceClient(conn)

	// Contact the server and print out its response.

	r, err := c.Hello(ctx, &helloworld.HelloMessage{
		Message: defaultName,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetResult())

	return nil

}
