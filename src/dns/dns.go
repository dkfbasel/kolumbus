package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/docker/docker/api/types"
	docker "github.com/moby/moby/client"
	"github.com/pkg/errors"
)

// DNS provides the methods
type DNS struct {
	DockerCli *docker.Client
	Services  map[string][]Endpoint
	sync.RWMutex
}

// Endpoint contains information for an endpoint instance of a given service
type Endpoint struct {
	Host string // host that the service is available on
	Port string // port of the service
}

// NewDNS will initialize a new dns server
func NewDNS() *DNS {

	dns := DNS{}
	dns.Services = make(map[string][]Endpoint)

	return &dns
}

// InitDocker will initialize a new docker command line interface
func (dns *DNS) InitDocker() error {

	cli, err := docker.NewEnvClient()
	if err != nil {
		return errors.Wrap(err, "could not initialize docker cli")
	}

	dns.DockerCli = cli

	return nil
}

// ServeHTTP will start an http server for administrative purposes
func (dns *DNS) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

// StartServer will start a http server
func (dns *DNS) StartServer() error {

	server := http.Server{
		Addr:    ":80",
		Handler: dns,
	}

	// start server and wait for completion
	return server.ListenAndServe()
}
