package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// HandleEnvoyRequest will handle requests to from envoyproxies to the EDS
// service (Endpoint-Discovery-Service)
func HandleEnvoyRequest(dns *DNS) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		response := EnvoyResponse{
			VersionInfo: "0",
		}

		response.Resources = []Resource{Resource{}}
		response.Resources[0].Type = "type.googleapis.com/envoy.api.v2.ClusterLoadAssignment"
		response.Resources[0].ClusterName = "helloworld_service"
		response.Resources[0].Endpoints = []Endpoint1{Endpoint1{}}

		response.Resources[0].Endpoints[0].LbEndpoints = []LbEndpoint{LbEndpoint{}}
		response.Resources[0].Endpoints[0].LbEndpoints[0].Endpoint = Endpoint2{}

		response.Resources[0].Endpoints[0].LbEndpoints[0].Endpoint.Address.SocketAddress.Address = "172.20.0.3"

		response.Resources[0].Endpoints[0].LbEndpoints[0].Endpoint.Address.SocketAddress.PortValue = "9211"

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

// EnvoyResponse is used to return information to envoy
type EnvoyResponse struct {
	VersionInfo string `json:"version_info"`
	TypeURL     string `json:"type_url,omitempty"`

	Resources []Resource `json:"resources"`
}

// Resource ..
type Resource struct {
	Type        string      `json:"@type"`
	ClusterName string      `json:"cluster_name"`
	Endpoints   []Endpoint1 `json:"endpoints"`
}

// Endpoint1 ..
type Endpoint1 struct {
	LbEndpoints []LbEndpoint `json:"lb_endpoints"`
}

// LbEndpoint ..
type LbEndpoint struct {
	Endpoint Endpoint2 `json:"endpoint"`
}

// Endpoint2 ..
type Endpoint2 struct {
	Address Address `json:"address"`
}

// Address ..
type Address struct {
	SocketAddress SocketAddress `json:"socket_address"`
}

// SocketAddress ..
type SocketAddress struct {
	Address   string `json:"address"`
	PortValue string `json:"port_value"`
}
