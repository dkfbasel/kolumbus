package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// HandleEnvoyClusterRequest will handle envoy discovery requests for clusters
func HandleEnvoyClusterRequest(dns *Kolumbus, errs chan<- error) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		response := EnvoyClusterResponse{
			VersionInfo: "1",
			TypeURL:     "type.googleapis.com/envoy.api.v2.Cluster",
		}

		response.Resources = []ClusterResource{}

		// lock write access to the dns struct
		dns.RLock()

		// iterate through all services
		for serviceName, hosts := range dns.Services {

			// initialize a new service cluster
			endpoint := ClusterResource{}

			endpoint.ResourceType = "type.googleapis.com/envoy.api.v2.Cluster"
			endpoint.Type = "strict_dns"

			// set the name of the service cluster
			endpoint.Name = fmt.Sprintf("%s_service_cluster", serviceName)
			endpoint.ConnectTimeout = "0.25s"

			// add all available hosts to the service cluster
			endpoint.Hosts = make([]Host, len(hosts))

			for index, item := range hosts {
				host := Host{
					SocketAddress: SocketAddress{
						Address: item.Host,
						Port:    item.Port,
					},
				}
				endpoint.Hosts[index] = host
			}

			response.Resources = append(response.Resources, endpoint)
		}

		dns.RUnlock()

		content, err := json.Marshal(response)
		if err != nil {
			errs <- errors.Wrap(err, "could not marshal response")
			return
		}

		_, err = w.Write(content)
		if err != nil {
			errs <- errors.Wrap(err, "could not write response")
		}

	}

}

// EnvoyClusterResponse is used to return information to envoy
type EnvoyClusterResponse struct {
	VersionInfo string            `json:"version_info"`
	Resources   []ClusterResource `json:"resources"`
	TypeURL     string            `json:"type_url"`
}

// ClusterResource ..
type ClusterResource struct {
	ResourceType   string   `json:"@type"`
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	ConnectTimeout string   `json:"connect_timeout"`
	LbPolicy       string   `json:"lb_policy,omitempty"`
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
