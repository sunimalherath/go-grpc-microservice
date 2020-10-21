package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/sunimalherath/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	// 1. Create a connection to the server - cc (ClientConnection)
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%v", err)
	}

	defer cc.Close()
	//c := sumpb.NewSumServiceClient(cc)
	c := calculatorpb.NewCalculatorServiceClient(cc)

	//doUnary(c)
	doServerStreaming(c)
}

// func doUnary(c sumpb.CalculatorServiceClient) {
// 	req := &sumpb.SumRequest{
// 		Sum: &sumpb.Sum{
// 			FirstNumber:  32,
// 			SecondNumber: 48,
// 		},
// 	}

// 	res, err := c.GetSum(context.Background(), req)
// 	if err != nil {
// 		log.Fatalf("Error while calling GetSum RPC: %v", err)
// 	}

// 	log.Printf("The result of Sum is : %v", res.Result)
// }

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a PrimeDecompositionServer streaming RPC...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecomposition RPC: %v", err)
	}

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving from server: %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}
