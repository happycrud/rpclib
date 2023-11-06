package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/happycurd/rpclib/protojson"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "127.0.0.1:9000", "the address to connect to")
	path = flag.String("path", "/helloworld.Greeter/SayHello", "grpc method path")
	data = flag.String("data", "{}", "grpc request data ")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.CallContentSubtype(protojson.JSON{}.Name())))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	resp := &protojson.Response{}
	if err := conn.Invoke(ctx, *path, *data, resp); err != nil {
		panic(err)
	}
	log.Printf("response: %+v", resp.Data)
}
