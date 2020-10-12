package main

import (
	"context"
	"log"

	"github.com/sunimalherath/grpc-go/calculator/sumpb"
	"google.golang.org/grpc"
)

func main() {
	// 1. Create a connection to the server - cc (ClientConnection)
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%v", err)
	}

	defer cc.Close()
	c := sumpb.NewSumServiceClient(cc)

	req := &sumpb.SumRequest{
		Sum: &sumpb.Sum{
			FirstNumber:  32,
			SecondNumber: 48,
		},
	}

	res, err := c.GetSum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GetSum RPC: %v", err)
	}

	log.Printf("The result of Sum is : %v", res.Result)
}
