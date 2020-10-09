package main

import (
	"fmt"
	"github.com/sunimalherath/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	// 1. Create a connection to server
	cc, err := grpc.Dial("localhost:500051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	// 3. Close the connection once all the code in main function are done
	defer cc.Close()

	// 2. Crete the client
	c := greetpb.NewGreetServiceClient(cc)
	fmt.Printf("Client created: %f", c)
}
