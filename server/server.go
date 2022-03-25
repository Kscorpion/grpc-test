package main

import (
	"GrpcTest/example"
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/golang/protobuf/proto"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	example.GoSpiderServer
}

//普通rpc调用
func (s *server) GetAddressResponse0(ctx context.Context, param *example.SendAddress) (*example.GetResponse, error) {
	log.Printf("Received: %v", param.GetAddress())
	return &example.GetResponse{
		HttpCode: proto.Int32(200),
		Response: proto.String("success"),
	}, nil
}

//客户端流式
func (s *server) GetAddressResponse1(srv example.GoSpider_GetAddressResponse1Server) error {
	for {
		//从流中获取消息
		res, err := srv.Recv()
		if err == io.EOF {
			//发送结果，并关闭
			return srv.SendAndClose(&example.GetResponse{
				HttpCode: proto.Int32(200),
				Response: proto.String("成功返回"),
			})
		}
		if err != nil {
			return err
		}
		log.Println(res.GetAddress(), res.GetMethod())
	}
}

//服务端流式
func (s *server) GetAddressResponse2(req *example.SendAddress, srv example.GoSpider_GetAddressResponse2Server) error {
	n := 1
	for {
		err := srv.Send(&example.GetResponse{HttpCode: proto.Int32(200), Response: proto.String(fmt.Sprintf("Address: %v, Method: %v", req.GetAddress(), req.GetMethod()))})
		if err != nil {
			return err
		}
		n++
		log.Printf("from stream client question:%d %s", n, req.GetAddress())
	}
}

//客户端服务端双流
func (s *server) GetAddressResponse3(srv example.GoSpider_GetAddressResponse3Server) error {
	n := 1
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = srv.Send(&example.GetResponse{HttpCode: proto.Int32(200), Response: proto.String("成功返回")})
		if err != nil {
			return err
		}
		n++
		log.Printf("from stream client question:%d %s", n, req.GetAddress())
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	example.RegisterGoSpiderServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
