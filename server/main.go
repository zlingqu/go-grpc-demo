package main

import (
	"context"
	"io"
	"log"
	"net"

	demo "go-grpc-demo/data"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type demoServer struct {
	demo.UnimplementedDemoServer
	savedResults []*demo.Response //用于服务端流
}

// 实现方法Add
func (s *demoServer) Add(ctx context.Context, in *demo.TwoNum) (*demo.Response, error) {
	x := in.X
	y := in.Y

	cRes := &demo.Response{
		Result: x + y,
	}

	return cRes, nil
}

// 实现方法SayHello
func (s *demoServer) SayHello(ctx context.Context, in *demo.HelloRequest) (*demo.HelloReply, error) {
	return &demo.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// 实现方法GetStream
func (s *demoServer) GetStream(in *demo.TwoNum, pipe demo.Demo_GetStreamServer) error {

	err := pipe.Send(&demo.Response{Result: in.X + in.Y}) //返回和
	if err != nil {
		return err
	}
	err = pipe.Send(&demo.Response{Result: in.X * in.Y}) //返回积
	if err != nil {
		return err
	}
	err = pipe.Send(&demo.Response{Result: in.X - in.Y}) //返回差
	if err != nil {
		return err
	}

	return nil
}

// 实现方法PutStream
func (s *demoServer) PutStream(pipe demo.Demo_PutStreamServer) error {
	var res int32
	for { //循环接收
		request, err := pipe.Recv()
		if err == io.EOF { //判断是否发送结束
			break
		}
		if err != nil {
			log.Println(err.Error())
		}
		res += request.X //累加
	}
	_ = pipe.SendAndClose(&demo.Response{Result: res}) //返回
	return nil
}

// 实现方法DoubleStream
func (s *demoServer) DoubleStream(pipe demo.Demo_DoubleStreamServer) error {

	for {
		request, err := pipe.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err = pipe.Send(&demo.Response{Result: request.X + request.Y}); err != nil {
			return err
		}

	}

}

func main() {
	const addr = ":50054"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	s := grpc.NewServer()
	demo.RegisterDemoServer(s, &demoServer{})
	reflection.Register(s)
	log.Printf("Server listeing at %s\n", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
