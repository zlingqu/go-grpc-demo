package main

import (
	"context"
	"flag"
	"fmt"
	demo "go-grpc-demo/data"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func printHelloResult(client demo.DemoClient, helloRequest *demo.HelloRequest) {
	fmt.Printf("请求参数是:name=%v\n", helloRequest.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.SayHello(ctx, helloRequest)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	fmt.Printf("返回内容:%v", resp.Message)
}

func printAddResult(client demo.DemoClient, twoNum *demo.TwoNum) {
	fmt.Printf("普通模式，请求参数是x=%d,y=%d\n", twoNum.X, twoNum.X)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.Add(ctx, twoNum)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	fmt.Printf("返回内容: %v", resp.Result)
}

func printGetStream(client demo.DemoClient, twoNum *demo.TwoNum) {
	fmt.Printf("请求参数是:x=%d,y=%d\n", twoNum.X, twoNum.Y)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	getStreamClient, err := client.GetStream(ctx, twoNum)
	if err != nil {
		log.Fatalf("error: %v: ", err)
	}
	for {
		resp, err := getStreamClient.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetStream错误=%v", client, err)
		}
		fmt.Printf("本次返回结果:%v\n", resp.Result)
	}
}

func printPutStream(client demo.DemoClient, oneNum []demo.OneNum) {
	// fmt.Printf("请求参数是:x=%v)", oneNum.X)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	putStreamClient, err := client.PutStream(ctx)
	if err != nil {
		log.Fatalf("error: %v: ", err)
	}
	for _, num := range oneNum {
		if err := putStreamClient.Send(&num); err != nil {

			log.Fatalf("PutStream错误=%v", err)
		}
	}
	resp, err := putStreamClient.CloseAndRecv()
	if err != nil {
		log.Fatalf("接收错误%v", err)
	}
	fmt.Printf("本次返回结果:%v\n", resp.Result)
}

func printDoubleStream(client demo.DemoClient) {
	// fmt.Printf("请求参数是:x=%v)", oneNum.X)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	doubleStreamClient, err := client.DoubleStream(ctx)
	if err != nil {
		log.Fatalf("error: %v: ", err)
	}
	ch := make(chan int32, 10)
	go func() { //并发编程，持续接收数据
		for {
			resp, err := doubleStreamClient.Recv()
			if err == io.EOF {
				close(ch)
				break
			}
			if err != nil {
				log.Fatalf("错误=%v", err)
			}
			ch <- resp.Result
		}

	}()

	go func() { //并发编程，持续发送10次请求
		for i := 0; i < 10; i++ {
			if err := doubleStreamClient.Send(&demo.TwoNum{X: int32(i), Y: int32(i)}); err != nil {
				log.Fatalf("发送数据失败: %v", err)
			} else {
				fmt.Println("发送数据", i, i)
			}

		}

	}()
	for k := 0; k < 10; k++ {
		fmt.Println("双向流： ", <-ch)
	}
	doubleStreamClient.CloseSend()
	<-ch

}

var (
	tls                = flag.Bool("tls", false, "是否使用tls")
	serverAddr         = flag.String("server_addr", "localhost:50054", "服务端地址，格式： host:port")
	serverHostOverride = flag.String("server_host_override", "a.grpc.test.com", "验证TLS握手返回的主机名的服务器名称。需要和服务端证书中dns段落匹配")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption

	if *tls {
		creds, err := credentials.NewClientTLSFromFile("keys/ca.crt", *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := demo.NewDemoClient(conn)

	fmt.Printf("#############第1次请求，简单模式########\n")
	printAddResult(client, &demo.TwoNum{X: 10, Y: 2})
	fmt.Printf("\n\n")

	fmt.Printf("#############第2次请求，简单模式########\n")
	printHelloResult(client, &demo.HelloRequest{Name: "张三"})
	fmt.Printf("\n\n")

	fmt.Printf("#############第3次请求，服务端流模式########\n")
	printGetStream(client, &demo.TwoNum{X: 10, Y: 2})
	fmt.Printf("\n\n")

	fmt.Printf("#############第4次请求，客户端流模式########\n")
	res := []demo.OneNum{
		demo.OneNum{
			X: 1,
		},
		demo.OneNum{
			X: 2,
		},
		demo.OneNum{
			X: 3,
		},
		demo.OneNum{
			X: 4,
		},
	}
	printPutStream(client, res)
	fmt.Printf("\n\n")

	fmt.Printf("#############第5次请求，双向流模式########\n")
	printDoubleStream(client)

}
