package main

import (
	"context"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/sunimalherath/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	result := "Hello "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// streaming ends
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error occured while client streaming: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += firstName + "! "
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		greetMessage := "Hello " + req.GetGreeting().GetFirstName() + " !"

		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: greetMessage,
		})
		if sendErr != nil {
			return sendErr
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051") // port binding
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()                            // create grpc server
	greetpb.RegisterGreetServiceServer(s, &server{}) // register server

	if err := s.Serve(lis); err != nil { // bind the server with the port
		log.Fatalf("Failed to serve: %v", err)
	}
}
