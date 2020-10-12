package main

import (
	"context"
	"log"
	"net"

	"github.com/sunimalherath/grpc-go/calculator/sumpb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Error creating listener: %v", err)
	}

	s := grpc.NewServer()
	sumpb.RegisterSumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}

func (*server) GetSum(ctx context.Context, req *sumpb.SumRequest) (*sumpb.SumResponse, error) {
	n1 := req.GetSum().FirstNumber
	n2 := req.GetSum().SecondNumber

	sum := n1 + n2

	res := &sumpb.SumResponse{
		Result: sum,
	}

	return res, nil
}
