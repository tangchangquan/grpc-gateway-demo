# gRPC Gateway

## 快速使用

- 启动项目

```shell
git clone https://github.com/helloworlde/grpc-gateway.git & cd grpc-gateway
make all 
```

- 访问

```shell
curl localhost:8090/hello\?message=world

{"result":"Hello world"}%
```

## 使用

### 安装依赖

- 安装 buf

buf 用于代替 protoc 进行生成代码，可以避免使用复杂的 protoc 命令，避免 protoc 各种失败问题

```shell
brew tap bufbuild/buf
brew install buf
```

- 安装 grpc-gateway

```shell
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
```

### 修改代码

- 添加 google.api 的 proto

添加 [`annotations.proto`](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/third_party/googleapis/google/api/annotations.proto)
和 [
http.proto`](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/third_party/googleapis/google/api/http.proto) 文件到 `
proto/google/api/`下；这两个文件用于支持 gRPC Gateway 代理

- 修改业务的 proto 文件

```diff
syntax = "proto3";

package io.github.helloworlde;
option go_package = "github.com/helloworlde/grpc-gateway;grpc_gateway";
option java_package = "io.github.helloworlde";
option java_multiple_files = true;
option java_outer_classname = "HelloGrpc";

+import "google/api/annotations.proto";

service HelloService{
  rpc Hello(HelloMessage) returns (HelloResponse){
+    option (google.api.http) = {
+      get: "/hello"
+    };
  }
}

message HelloMessage {
  string message = 1;
}

message HelloResponse {
  string result = 1;
}
```

### 配置 Gateway

- 添加 buf 配置文件 buf.gen.yaml

```yaml
version: v1beta1
plugins:
  - name: go
    out: proto
    opt: paths=source_relative
  - name: go-grpc
    out: proto
    opt: paths=source_relative,require_unimplemented_servers=false
  - name: grpc-gateway
    out: proto
    opt: paths=source_relative
```

- 添加配置文件 buf.yaml

```yaml
version: v1beta1
build:
  roots:
    - proto
```

- 生成 Gateway 的代码

会生成 `*.gw.go` 格式的文件，该文件是 gRPC Gateway 代理具体服务的实现

```shell
buf generete
```

- 添加 gRPC Gateway Proxy

```go
func StartGwServer() {
	conn, _ := grpc.DialContext(
		context.Background(),
		"0.0.0.0:9090",
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)

	mux := runtime.NewServeMux()
	// 注册服务
	pb.RegisterHelloServiceHandler(context.Background(), mux, conn)

	server := &http.Server{
		Addr:    ":8090",
		Handler: mux,
	}

	server.ListenAndServe()
}
```

启动应用后即可访问相应的接口

## 参考文档

- [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
- [buf](https://buf.build/)
- [grpc-gateway document](https://grpc-ecosystem.github.io/grpc-gateway/)