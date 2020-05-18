package main

import (
	"context"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/moby/moby/client"
	"github.com/pkg/errors"
)

// StartDockerWatch will initialize a new docker command line interface and
// start watching for container changes in docker
func (dns *Kolumbus) StartDockerWatch(config Config, errs chan<- error) {

	cli, err := docker.NewEnvClient()
	if err != nil {
		errs <- errors.Wrap(err, "could not initialize docker cli")
		return
	}

	// initialize periodic watcher for clients
	go func() {

		for {

			// reload container info to check if new network informations are available
			err = dns.loadContainerInfo(cli, config.Hostname)
			if err != nil {
				errs <- errors.Wrap(err, "could not reload container network informations")
			}

			// find all docker services
			services, err := findServices(cli, config, dns.getContainerNetworkIDs())

			if err != nil {
				errs <- errors.Wrap(err, "could not find docker services")
			}

			if err == nil {
				// replace the information in the dns store
				dns.Lock()
				dns.Services = services
				dns.Unlock()
			}

			// wait 5 seconds to check again
			<-time.After(config.DiscoveryInterval)
		}
	}()
}

// getContainerNetworkIDs will return a list of network ids where the kolumbus
// docker container is a member
func (dns *Kolumbus) getContainerNetworkIDs() []string {

	networkIDs := []string{}

	// return an empty list if container information are not available
	if dns.ContainerInfo == nil || dns.ContainerInfo.NetworkSettings == nil {
		return networkIDs
	}

	// get all network ids of the kolumbus docker container
	for _, network := range dns.ContainerInfo.NetworkSettings.Networks {
		networkIDs = append(networkIDs, network.NetworkID)
	}

	return networkIDs
}

// findServices will handle the recovery of services via docker labels
func findServices(cli *docker.Client, config Config, networks []string) (map[string][]Endpoint, error) {

	// setup a new container list filter for containers that share a network
	// with the kolumbus docker container
	containerFilters := filters.NewArgs()
	for _, networkID := range networks {
		containerFilters.Add("network", networkID)
	}

	// get a list of running containers filtered by networks through the docker cli
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All:     false,
		Filters: containerFilters,
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch the list of docker containers")
	}

	// clear the list of services that are available
	services := make(map[string][]Endpoint)

	for _, container := range containers {

		// ignore all containers that do not have the label envoyproxy.service
		if _, ok := container.Labels["envoyproxy.service"]; !ok {
			continue
		}

		// get the name of the service that should be available for other
		// envoy proxies
		name := container.Labels["envoyproxy.service"]

		// get the host id and port (defaults to port 80)
		// the short code of the host id is used as alias for networking
		host := shortID(container.ID, 12)

		// get the envoyproxy port from labels (default is 80)
		port, ok := container.Labels["envoyproxy.port"]
		if !ok {
			port = "80"
		}

		// create an endpoint
		endpoint := Endpoint{
			Host: host,
			Port: port,
		}

		// add the container as endpoint for the given service
		endpoints := services[name]
		endpoints = append(endpoints, endpoint)
		services[name] = endpoints
	}

	log.Printf("envoy services: %+v\n", services)

	return services, nil
}

// shortID will take the first x characters from the given id
func shortID(containerID string, x int) string {
	runes := []rune(containerID)
	return string(runes[:x])
}
