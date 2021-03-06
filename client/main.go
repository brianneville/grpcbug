package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/brianneville/grpcbug/defaults"
	"github.com/brianneville/grpcbug/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "", "address of server")
	var useWorkaround bool
	flag.BoolVar(&useWorkaround, "work", false, "use client-side workaround")
	flag.Parse()
	if addr == "" {
		flag.Usage()
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if useWorkaround {
		cfg = append(defaults.WorkaroundCfg(), cfg...)
	}

	cc, err := grpc.DialContext(ctx, addr+":"+defaults.DefaultPort,
		cfg...)
	if err != nil {
		log.Fatal(err)
	}
	client := proto.NewMockRPCClient(cc)
	resp, err := client.Get(ctx, &proto.GetRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("recieved response, length=%d\n", len(resp.B))
}
