package main

import (
	"github.com/sunimalherath/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {

}

func main() {
	lis, err := net.Listen("tcp","0.0.0.0:50051") // port binding
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer() // create grpc server
	greetpb.RegisterGreetServiceServer(s,&server{}) // register server

	if err := s.Serve(lis); err != nil { // bind the server with the port
		log.Fatalf("Failed to serve: %v", err)
	}
}
