package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	helloworld "github.com/dkfbasel/kolumbus/examples/grpc/src/proto"
	"google.golang.org/grpc"
)

func main() {

	var port int
	flag.IntVar(&port, "port", 8080, "grpc port")
	flag.Parse()

	// initialize a tcp listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	// create a new grpc server that is secured with ssl
	srv := grpc.NewServer()

	// register our custom server
	helloworld.RegisterHelloWorldServer(srv, newHelloWorldServer())

	log.Printf("starting grpc server on port %d\n", port)

	// start the server
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}

}

// server implements the service interface
type server struct{}

func newHelloWorldServer() *server {
	return &server{}
}

// Echo will return the message as echo
func (s *server) Echo(ctx context.Context, in *helloworld.EchoRequest) (*helloworld.EchoResponse, error) {
	log.Printf("- echoing: %s", in.Message)

	return &helloworld.EchoResponse{
		Message: fmt.Sprintf("echo: %s", in.Message),
	}, nil
}
