package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// HandleEnvoyClusterRequest will handle envoy discovery requests for clusters
func HandleEnvoyClusterRequest(dns *DNS) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		log.Println("got cluster request")

		response := EnvoyClusterResponse{
			VersionInfo: "0",
		}

		endpoint := ClusterResource{}

		endpoint.ResourceType = "type.googleapis.com/envoy.api.v2.Cluster"
		endpoint.Type = "strict_dns"
		endpoint.Name = "helloworld_service"
		endpoint.ConnectTimeout = "0.25s"

		endpoint.Hosts = []Host{Host{
			SocketAddress: SocketAddress{
				Address: "helloworld_service",
				Port:    "9211",
			},
		}}

		response.Resources = []ClusterResource{endpoint}

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

// EnvoyClusterResponse is used to return information to envoy
type EnvoyClusterResponse struct {
	VersionInfo string            `json:"version_info"`
	Resources   []ClusterResource `json:"resources"`
}

// ClusterResource ..
type ClusterResource struct {
	ResourceType   string   `json:"@type"`
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	ConnectTimeout string   `json:"connect_timeout"`
	HTTP2          struct{} `json:"http2_protocol_options"`
	Hosts          []Host   `json:"hosts"`
}

// Host ..
type Host struct {
	SocketAddress SocketAddress `json:"socket_address"`
}

// SocketAddress ..
type SocketAddress struct {
	Address string `json:"address"`
	Port    string `json:"port_value"`
}
