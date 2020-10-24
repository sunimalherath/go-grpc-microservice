package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	//doServerStreaming(c)
	//doAverage(c)
	doErrorUnary(c)
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

func doAverage(c calculatorpb.CalculatorServiceClient) {
	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 23,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 45,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 74,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 28,
		},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error occured while steaming: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req.GetNumber())
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	// Start alternate method ->  or can do as follows with a number slice
	numbers := []int64{23, 45, 74, 28}

	for _, number := range numbers {
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}
	// End alternate method

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving results: %v", err)
	}
	fmt.Printf("Computed average: %v", res.GetAverage())
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	// correct call
	doErrorCall(c, 20)

	// incorrect call
	doErrorCall(c, -39)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// this is the actual error from gRPC
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			// can check for different errors like below
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("Sent a negative number!")
			}
		} else {
			log.Fatalf("Some other error not in status: %v", err)
		}
	}
	fmt.Printf("Result of square root of %v: %v\n", n, res.GetNumberRoot())
}
