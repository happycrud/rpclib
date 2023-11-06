package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/happycurd/rpclib/discovery"
	pb "github.com/happycurd/rpclib/example/helloworld/helloworld"
	_ "github.com/happycurd/rpclib/protojsonserver"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	discovery.InitEtcdClient("http://127.0.0.1:2379")
	si := discovery.NewServiceInstance("helloworld", "172.16.109.105:"+fmt.Sprintf("%d", *port), discovery.GetEtcdClient())
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	err = si.Register()
	if err != nil {
		panic(err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
