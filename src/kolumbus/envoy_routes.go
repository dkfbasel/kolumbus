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

		// initialize all resources
		response.Resources = []RouteResource{}

		// lock access to dns records
		dns.RLock()

		// iterate through all services
		for serviceName := range dns.Services {

			// initalize a new endpoint
			endpoint := RouteResource{}

			endpoint.ResourceType = "type.googleapis.com/envoy.api.v2.RouteConfiguration"
			endpoint.Name = "kolumbus_routes"

			// define the route that will be matched for the service
			routeMatch := fmt.Sprintf("/%s", serviceName)

			// define the name of the cluster for the service
			clusterName := fmt.Sprintf("%s_service_cluster", serviceName)

			// generate the virtual hosts configuration for the service
			endpoint.VirtualHosts = []VirtualHost{VirtualHost{
				Name:    fmt.Sprintf("%s-virtual-hosts", serviceName),
				Domains: []string{"*"},
				Routes: []Route{Route{
					Match: RouteMatch{
						Prefix: routeMatch,
					},
					Route: RouteRouting{
						Cluster: clusterName,
					},
				},
				},
			}}

			// append the service cluster to the routes list
			response.Resources = append(response.Resources, endpoint)

		}

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
