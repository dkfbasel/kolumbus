package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/pkg/errors"
)

// Kolumbus provides the methods
type Kolumbus struct {
	Services map[string][]Endpoint
	sync.RWMutex
}

// Endpoint contains information for an endpoint instance of a given service
type Endpoint struct {
	Host string // host that the service is available on
	Port string // port of the service
}

// FindABraveNewWorld will initialize a new kolumbus dns server
func FindABraveNewWorld() *Kolumbus {

	dns := Kolumbus{}
	dns.Services = make(map[string][]Endpoint)

	return &dns
}

// StartEnvoyDataPlaneServer will start a http server to provide service
// information for envoy proxies
func (dns *Kolumbus) StartEnvoyDataPlaneServer(address int, errs chan<- error) {

	mux := http.NewServeMux()
	mux.HandleFunc("/v2/discovery:clusters", HandleEnvoyClusterRequest(dns, errs))
	mux.HandleFunc("/v2/discovery:routes", HandleEnvoyRouteRequest(dns, errs))
	mux.HandleFunc("/v1/certs/list/approved", HandleEnvoyCertificateRequest(dns, errs))
	mux.HandleFunc("/", HandleAnyRequest(dns, errs))

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", address), mux)
		if err != nil {
			errs <- errors.Wrap(err, "could not start discovery service")
		}
	}()
}

// HandleAnyRequest will handle all request to paths that were not specified before
func HandleAnyRequest(dns *Kolumbus, errs chan<- error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("got other request: %s\n", r.URL.String())

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errs <- errors.Wrap(err, "could not read body")
		}
		_ = r.Body.Close()

		log.Printf("body: %s\n", body)
	}
}
