package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/happycurd/rpclib/discovery"
	pb "github.com/happycurd/rpclib/example/helloworld/helloworld"
	"github.com/happycurd/rpclib/protojson"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()

	conn, _ := discovery.NewConn("helloworld")
	// Set up a connection to the server.
	// conn, err := grpc.Dial("etcd:///"+"helloworld", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.CallContentSubtype(protojson.JSON{}.Name())))
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	for i := 0; i < 10; i++ {
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "xx"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.GetMessage())
	}

	req := `{"name":"world"}`
	resp := &protojson.Response{}
	err := conn.Invoke(ctx, "/helloworld.Greeter/SayHello", req, resp)
	log.Printf("Greeting: %+v err:%v", resp, err)
}
