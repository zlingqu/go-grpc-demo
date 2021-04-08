# 1. 基于go的grpc的一个demo

有关grpc协议，可以google搜索下，可以和http协议做对比理解和学习。k8s中大量使用此种协议。

# 2. 目录说明
## 2.1 data
注意：这部分代码，对于2.2、2.3是公用的

### add.proto
protobuf 即 Protocol Buffers，是一种轻便高效的结构化数据存储格式。可以理解为和json、XML类似和语言无关，对别json、xml等有其自身的优势，详情可google下。


grpc中使用此种方案对数据进行序列化、反序列化。

add.proto即是protobuf的源文件，该文件可以借助protoc工具(https://github.com/protocolbuffers/protobuf/releases)将其格式化成各种语言的代码，比如go、java、python、js、php等，详情可搜索下相关使用方法。
### add.pb.go
将add.proto转换成的go代码
```bash
# 安装插件
go get -u github.com/golang/protobuf/protoc-gen-go

# 生成go代码
protoc --go_out=.  add.proto 
# go_out表示格式化成go代码，类似的还有java_out、python_out等
# go_out=.，最后的.表示生成的文件在当前目录
# add.proto，表示执行的目标文件，这里使用相对路径，也可以用类似于 /root/abc/*.proto等方式
# 生成的文件名是固定的, 文件名+.pb+.go
```
### add_grpc.pb.go


将add.proto转换成的go grpc需要的相关代码
```bash
# 安装插件
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
# 生成go grpc相关的代码
protoc --go-grpc_out=. add.proto
# 可以和上一步写到一起 protoc --go_out=. --go-grpc_out=. add.proto 
# --go-grpc_out需要安装好插件 
```

## 2.2 server

服务端代码，运行方式
```bash
go run server/main.go
```
## 2.3 client

客户端代码，运行方式
```bash
go run client/main.go
```
另外也有类似于grpcui的可视化客户端

# 3 可视化调试工具

类似于http的调试工具postman一样，grpc也有可视化工具，见：https://github.com/fullstorydev/grpcui

可以直接下载二进制文件后，运行


```
grpcui -plaintext 127.0.0.1:50052
# 127.0.0.1:50052 表示grpc服务的ip和端口
```