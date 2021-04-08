package main

import (
	"context"
	"fmt"
	calcpb "go-grpc-demo/data"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:50054", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	addServiceClient := calcpb.NewCalculatorClient(conn)
	response, err := addServiceClient.Calculate(context.Background(), &calcpb.CalculatorRequest{X: 1, Y: 20})
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Result)
}
