package main

import (
	"context"
	"log"
	"net"

	calcpb "go-grpc-demo/data"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type service struct {
	calcpb.UnimplementedCalulatorServer
}

func (cs *service) Calculate(ctx context.Context, cReq *calcpb.CalulatorRequest) (*calcpb.CalulatorResponse, error) {
	x := cReq.X
	y := cReq.Y

	log.Printf("Received request to calculate %d + %d\n", x, y)
	cRes := &calcpb.CalulatorResponse{
		Result: x + y,
	}

	return cRes, nil
}

func main() {
	const addr = ":50054"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	s := grpc.NewServer()
	calcpb.RegisterCalulatorServer(s, &service{})
	reflection.Register(s)
	log.Printf("Server listeing at %s\n", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
