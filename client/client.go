package main

import (
	"GrpcTest/example"
	"context"
	"io"
	"log"

	"time"

	"github.com/golang/protobuf/proto"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

//普通rpc调用
func normalCall(ctx context.Context, client example.GoSpiderClient) {
	r, err := client.GetAddressResponse0(ctx, &example.SendAddress{Address: proto.String("家"), Method: proto.String("1801")})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("普通调用 --- Greeting: %s %d", r.GetResponse(), r.GetHttpCode())
}

//客户端流式调用
func clientStreamCall(ctx context.Context, times int, client example.GoSpiderClient) {
	//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	stream, err := client.GetAddressResponse1(ctx)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
	for n := 0; n < times; n++ {
		//向流中发送消息
		err := stream.Send(&example.SendAddress{
			Address: proto.String("home"),
			Method:  proto.String("客户端流式"),
		})
		if err != nil {
			log.Fatalf("stream request err: %v", err)
		}
	}
	//关闭流并获取返回的消息
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("RouteList get response err: %v", err)
	}
	log.Println(res)
}

//服务端流式调用
func serverStreamCall(ctx context.Context, client example.GoSpiderClient) {
	stream, err := client.GetAddressResponse2(ctx, &example.SendAddress{
		Address: proto.String("home"),
		Method:  proto.String("服务端流式"),
	})
	if err != nil {
		log.Fatalf("Call ListStr err: %v", err)
	}
	for {
		res, err := stream.Recv()
		// 判断消息流是否已经结束
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("ListStr get stream err: %v", err)
		}
		// 打印返回值
		log.Println(res.GetHttpCode(), res.GetResponse())
	}
}

//客户端服务端双流式调用
func clientStreamServerStreamCall(ctx context.Context, times int, client example.GoSpiderClient) {
	//请求返回均流式
	stream, err := client.GetAddressResponse3(ctx)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	for n := 0; n < times; n++ {
		err := stream.Send(&example.SendAddress{Address: proto.String("home"), Method: proto.String("双流式调用")})
		if err != nil {
			log.Fatalf("流 -- stream request err: %v", err)
		}
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Conversations get stream err: %v", err)
		}
		// 打印返回值
		log.Println(res.GetHttpCode(), res.GetResponse())
	}
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := example.NewGoSpiderClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	//normalCall(ctx, client)
	//clientStreamCall(ctx, 5, client)
	//serverStreamCall(ctx, client)
	clientStreamServerStreamCall(ctx, 5, client)
}
