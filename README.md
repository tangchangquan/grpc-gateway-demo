# gRPC Gateway

## 快速使用

- 启动项目

```shell
git clone https://github.com/helloworlde/grpc-gateway.git & cd grpc-gateway
go mod tidy 
go run main.go
```

- 访问

```shell
curl localhost:8090/hello\?message=world

{"result":"Hello world"}%
```