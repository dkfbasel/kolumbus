package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// HandleEnvoyRouteRequest will handle envoy discovery requests for routes
func HandleEnvoyRouteRequest(dns *DNS) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		log.Println("got routes request")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("could not read body: %+v\n", err)
		}
		_ = r.Body.Close()

		log.Printf("body: %s\n", body)

		endpoint := RouteResource{}

		endpoint.ResourceType = "type.googleapis.com/envoy.api.v2.RouteConfiguration"
		endpoint.Name = "helloworld_service_routes"

		endpoint.VirtualHosts = []VirtualHost{VirtualHost{
			Name:    "client-virtual-hosts",
			Domains: []string{"*"},
			Routes: []Route{Route{
				Match: RouteMatch{
					Prefix: "/",
				},
				Route: RouteRouting{
					Cluster: "helloworld_service_cluster",
				},
			},
			},
		}}

		response := EnvoyRouteResponse{
			VersionInfo: "1",
			TypeURL:     "type.googleapis.com/envoy.api.v2.RouteConfiguration",
		}

		response.Resources = []RouteResource{endpoint}

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
