package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"

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

// StartServer will start a http server
func (dns *DNS) StartServer() error {

	mux := http.NewServeMux()
	mux.HandleFunc("/dns", HandleServiceDiscovery(dns))
	mux.HandleFunc("/v2/discovery:clusters", HandleEnvoyClusterRequest(dns))
	mux.HandleFunc("/v2/discovery:routes", HandleEnvoyRouteRequest(dns))
	mux.HandleFunc("/", HandleAnyRequest(dns))

	// start server and wait for completion
	return http.ListenAndServe(":80", mux)
}

// HandleAnyRequest will handle all request to paths that were not specified before
func HandleAnyRequest(dns *DNS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("got other request: %s\n", r.URL.String())

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("could not read body: %+v\n", err)
		}
		_ = r.Body.Close()

		log.Printf("body: %s\n", body)
	}
}
