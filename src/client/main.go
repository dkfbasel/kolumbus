package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	helloworld "bitbucket.org/dkfbasel/envoy/src/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	var address string
	var port int
	var certificates string

	flag.IntVar(&port, "port", 8081, "web hosting port")
	flag.StringVar(&address, "grpc", "8082", "grpc service address")
	flag.StringVar(&certificates, "certificates", "/app/certificates", "certificate directory")
	flag.Parse()

	if strings.Contains(address, ":") == false {
		address = fmt.Sprintf("localhost:%s", address)
	}

	// load certificate file
	creds, err := credentials.NewClientTLSFromFile(filepath.Join(certificates, "service.cert"), "")
	if err != nil {
		log.Fatalf("coud not load certificate files: %v\n", err)
	}

	// set up a connection to the grpc server
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}

	defer conn.Close() // nolint: errcheck

	client := helloworld.NewHelloWorldClient(conn)

	log.Printf("starting client with connection to grpc: %s\n", address)

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		message := r.URL.Query().Get("message")

		if message == "" {
			fmt.Fprintln(w, "Message must be specified")
			return
		}

		ctx := context.Background()
		result, err2 := client.Echo(ctx, &helloworld.EchoRequest{
			Message: message,
		})

		if err2 != nil {
			log.Fatalf("could not echo: %v\n", err)
			return
		}

		log.Printf("- %s\n", result)
		fmt.Fprintln(w, result)
		return

	}

	log.Printf("starting html server on port %d\n", port)

	http.HandleFunc("/", handlerFunc)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("could not start server on port %d\n", port)
	}

}
