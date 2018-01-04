package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
)

// HandleServiceDiscovery will handle the recovery of services via docker
func HandleServiceDiscovery(dns *DNS) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		log.Println("got http request")

		containers, err := dns.DockerCli.ContainerList(context.Background(), types.ContainerListOptions{All: false})

		if err != nil {
			log.Fatalln(err)
		}

		// lock access to the dns service
		dns.Lock()

		// clear the list of services that are available
		dns.Services = make(map[string][]Endpoint)

		for _, container := range containers {

			if _, ok := container.Labels["envoyproxy.service"]; !ok {
				continue
			}

			// get the name of the service that envoy should proxy to
			name := container.Labels["envoyproxy.service"]

			// get the host name and port
			host := container.Names[0]
			port, ok := container.Labels["envoyproxy.port"]
			if !ok {
				port = "80"
			}

			// create an endpoint
			endpoint := Endpoint{
				Host: host,
				Port: port,
			}

			endpoints := dns.Services[name]
			endpoints = append(endpoints, endpoint)
			dns.Services[name] = endpoints

			// print the container name
			fmt.Fprintf(w, "# %+v\n", container.Names)
		}

		// unlock access
		dns.Unlock()

	}

}
