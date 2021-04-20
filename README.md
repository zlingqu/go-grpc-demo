# 1. 基于go的grpc的一个demo

有关grpc协议，可以google搜索下，可以和http协议做对比理解和学习。k8s中大量使用此种协议。

#2. 功能说明

该示例项目演示了grpc的四种服务方式，分别是

```go
// 简单模式。一个请求，一个响应。
rpc Add (TwoNum) returns (Response) {} //客户端发送一个请求，包含两个数字，服务端是返回两个数字的和
rpc SayHello (HelloRequest) returns (HelloReply) {} //发送一个name字符串，返回hello name

//服务端流模式，客户端发送一个请求，服务端返回多次。
rpc GetStream (TwoNum) returns (stream Response) {} //请求一次，返回三次，分别是两数子和、两数之积、两数之差

//客户端流模式，客户端发送多次请求，服务端响应一次。
rpc PutStream (stream OneNum) returns (Response) {}//请求中每次都是一个数字，发送完成后，服务端返回所有数字之和

//双向流，发送和接收同时进行，互不干扰
rpc DoubleStream (stream TwoNum) returns (stream Response) {} //每次请求都返回两个数字之和
```


# 2. 目录说明
## 2.1 data
注意：这部分代码，对于2.2、2.3是公用的

#### 2.1.1 add.proto
protobuf 即 Protocol Buffers，是一种轻便高效的结构化数据存储格式。可以理解为和json、XML类似和语言无关，对别json、xml等有其自身的优势，详情可google下。


grpc中使用此种方案对数据进行序列化、反序列化。

add.proto即是protobuf的源文件，该文件可以借助protoc工具(https://github.com/protocolbuffers/protobuf/releases) 将其格式化成各种语言的代码，比如go、java、python、js、php等，详情可搜索下相关使用方法。
#### 2.1.2 add.pb.go
将add.proto转换成的go代码
```bash
# 安装插件，用于--go_out参数
go get -u github.com/golang/protobuf/protoc-gen-go

# 生成go代码
protoc --go_out=.  add.proto 
# go_out表示格式化成go代码，类似的还有java_out、python_out等
# go_out=.，最后的.表示生成的文件在当前目录
# add.proto，表示执行的目标文件，这里使用相对路径，也可以用类似于 /root/abc/*.proto等方式
# 生成的文件名是固定的, 文件名+.pb+.go
```
#### 2.1.3 add_grpc.pb.go


将add.proto转换成的go grpc需要的相关代码
```bash
# 安装插件,用于--go-grpc_out参数
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
# 生成go grpc相关的代码
protoc --go-grpc_out=. add.proto
# 可以和上一步写到一起
protoc --go_out=. --go-grpc_out=. add.proto 
```

## 2.2 server

server目录，服务端代码，运行方式
```bash
go run server/main.go

#如果使用tls,使用
go run server/main.go -tls=true

```
## 2.3 client

client目录，客户端代码，运行方式
```bash
go run client/main.go

#如果使用tls,使用-tls参数
go run client/main.go -tls=true

#指定服务端地址，使用-server_addr参数
go run client/main.go -server_addr="localhost:50054"

```
输出内容类似如下
```bash

#############第1次请求，简单模式########
普通模式，请求参数是x=10,y=10
返回内容: 12

#############第2次请求，简单模式########
请求参数是:name=张三
返回内容:Hello 张三

#############第3次请求，服务端流模式########
请求参数是:x=10,y=2
本次返回结果:12
本次返回结果:20
本次返回结果:8


#############第4次请求，客户端流模式########
本次返回结果:10


#############第5次请求，双向流模式########
发送数据 0 0
发送数据 1 1
双向流：  0
发送数据 2 2
发送数据 3 3
发送数据 4 4
发送数据 5 5
发送数据 6 6
发送数据 7 7
双向流：  2
双向流：  4
双向流：  6
双向流：  8
双向流：  10
双向流：  12
双向流：  14
双向流：  16
发送数据 8 8
发送数据 9 9
双向流：  18
```


# 3 可视化调试工具

类似于http的调试工具postman一样，grpc也有可视化工具，见：https://github.com/fullstorydev/grpcui

可以直接下载二进制文件后，运行


```
grpcui -plaintext 127.0.0.1:50054
# 127.0.0.1:50054 表示grpc服务的ip和端口
```