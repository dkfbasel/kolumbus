package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
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

// StartDNSServer will start a http server to provide service information for envoy
// proxies
func (dns *Kolumbus) StartDNSServer() error {

	mux := http.NewServeMux()
	mux.HandleFunc("/v2/discovery:clusters", HandleEnvoyClusterRequest(dns))
	mux.HandleFunc("/v2/discovery:routes", HandleEnvoyRouteRequest(dns))
	mux.HandleFunc("/", HandleAnyRequest(dns))

	// start server and wait for completion
	return http.ListenAndServe(":80", mux)
}

// HandleAnyRequest will handle all request to paths that were not specified before
func HandleAnyRequest(dns *Kolumbus) http.HandlerFunc {
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
