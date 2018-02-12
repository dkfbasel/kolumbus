package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// HandleEnvoyRouteRequest will handle envoy discovery requests for routes
func HandleEnvoyRouteRequest(dns *Kolumbus) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// log.Println("got routes request")

		// body, err := ioutil.ReadAll(r.Body)
		// if err != nil {
		// 	log.Printf("could not read body: %+v\n", err)
		// }
		// _ = r.Body.Close()

		// log.Printf("body: %s\n", body)

		// initialize the envoy route reponse
		response := EnvoyRouteResponse{
			VersionInfo: "1",
			TypeURL:     "type.googleapis.com/envoy.api.v2.RouteConfiguration",
		}

		// lock access to dns records
		dns.RLock()

		// initalize a new endpoint
		endpoint := RouteResource{}

		endpoint.ResourceType = "type.googleapis.com/envoy.api.v2.RouteConfiguration"
		endpoint.Name = "kolumbus_routes"

		// define a new virtual host
		virtualHost := VirtualHost{
			Name:    "kolumbus-virtual-hosts",
			Domains: []string{"*"},
			Routes:  []Route{},
		}

		// iterate through all services
		for serviceName := range dns.Services {

			// define a new route to match (one for each service cluster)
			route := Route{
				Match: RouteMatch{
					// prefix must match the name of the grpc service
					Prefix: fmt.Sprintf("/%s", serviceName),
				},
				Route: RouteRouting{
					// cluster to route to
					Cluster: fmt.Sprintf("%s_service_cluster", serviceName),
				},
			}

			// append the service cluster to the routes list
			virtualHost.Routes = append(virtualHost.Routes, route)

		}

		// only one virtual host is used
		endpoint.VirtualHosts = []VirtualHost{virtualHost}

		// only one endpoint resources is provided (all through kolumbus)
		response.Resources = []RouteResource{endpoint}

		// unlock access
		dns.RUnlock()

		content, err := json.Marshal(response)
		if err != nil {
			log.Printf("could not marshal response: %+v\n", err)
		}

		_, err = w.Write(content)
		if err != nil {
			log.Printf("problem when writing: %+v\n", err)
		}

	}

}

// EnvoyRouteResponse is used to return information to envoy
type EnvoyRouteResponse struct {
	VersionInfo string          `json:"version_info"`
	TypeURL     string          `json:"type_url"`
	Resources   []RouteResource `json:"resources"`
}

// RouteResource ..
type RouteResource struct {
	ResourceType string        `json:"@type"`
	Name         string        `json:"name"`
	VirtualHosts []VirtualHost `json:"virtual_hosts"`
}

// VirtualHost ..
type VirtualHost struct {
	Name    string   `json:"name"`
	Domains []string `json:"domains"`
	Routes  []Route  `json:"routes"`
}

// Route ..
type Route struct {
	Match RouteMatch   `json:"match"`
	Route RouteRouting `json:"route"`
}

// RouteMatch ..
type RouteMatch struct {
	Prefix string `json:"prefix"`
}

// RouteRouting ..
type RouteRouting struct {
	Cluster string `json:"cluster"`
}
