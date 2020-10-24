package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/sunimalherath/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	// 1. Create a connection to server
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	// 3. Close the connection once all the code in main function are done
	defer cc.Close()

	// 2. Crete the client
	c := greetpb.NewGreetServiceClient(cc)

	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDiStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Doe",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Jane",
			LastName:  "Doe",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Erro while reading stream: %v", err)
		}

		log.Printf("Response from GreetingManyTimes: %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jane",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Robert",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Neo",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sneding request: %v", req.GetGreeting().GetFirstName())
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receiving response: ", err)
	}
	fmt.Printf("Long greet response: %v", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Robert",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jane",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Natasha",
			},
		},
	}

	waitc := make(chan struct{})
	// Sending stream of messages
	go func() {
		for _, req := range requests {
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("Error while sending data: %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// Receiving stream of messages
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("End receiving data")
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving data: %v", err)
				break
			}
			fmt.Printf("Received: %v \n", res.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}
