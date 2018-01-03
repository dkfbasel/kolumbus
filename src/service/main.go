package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"path/filepath"

	helloworld "bitbucket.org/dkfbasel/envoy/src/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	var port int
	var certificates string
	flag.IntVar(&port, "port", 8082, "grpc port")
	flag.StringVar(&certificates, "certificates", "/app/certificates", "certificate directory")

	flag.Parse()

	// sslCertificate, err := ioutil.ReadFile(filepath.Join(certificates, "service.cert"))
	// if err != nil {
	// 	log.Fatalf("could not read ssl certificate file")
	// }
	//
	// sslKey, err := ioutil.ReadFile(filepath.Join(certificates, "service.key"))
	// if err != nil {
	// 	log.Fatalf("could not read ssl key file")
	// }
	//
	// cert, err := tls.X509KeyPair(sslCertificate, sslKey)
	// if err != nil {
	// 	log.Fatalf("could not generate key pair")
	// }

	creds, err := credentials.NewServerTLSFromFile(filepath.Join(certificates, "service.cert"), filepath.Join(certificates, "service.key"))
	if err != nil {
		log.Fatalf("could not read certificate files: %v\n", err)
	}

	// initialize a tcp listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	// create a new grpc server that is secured with ssl
	srv := grpc.NewServer(grpc.Creds(creds))

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
