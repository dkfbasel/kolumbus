package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	helloworld "github.com/dkfbasel/kolumbus/examples/grpc-dev/src/proto/helloworld"
	"google.golang.org/grpc"
)

func main() {

	var address string
	var port int

	flag.IntVar(&port, "port", 80, "web hosting port")
	flag.StringVar(&address, "kolumbus", "kolumbus:1494", "kolumbus envoyproxy address")
	flag.Parse()

	// read kolumbus address from the environment
	envAddress := os.Getenv("kolumbus")
	if envAddress != "" {
		address = envAddress
	}

	if strings.Contains(address, ":") == false {
		log.Fatalf("host for grpc connection must be specified")
	}

	// set up a connection to the grpc server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
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
			log.Printf("could not call service: %v\n", err2)
			fmt.Fprintln(w, "could not call service")
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
