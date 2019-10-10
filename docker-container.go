package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	docker "github.com/moby/moby/client"
	"github.com/pkg/errors"
)

// IdentifyContainer will identify the kolumbus docker container
func (dns *Kolumbus) IdentifyContainer(config Config, errs chan<- error) {

	cli, err := docker.NewEnvClient()
	if err != nil {
		errs <- errors.Wrap(err, "could not initialize docker cli")
		return
	}

	done := make(chan bool)

	// initialize periodic watcher to identify the kolumbus docker container
	// until found or a timeout occures after 30 seconds
	go func() {

		for {
			// get the kolumbus docker container information by matching hostname
			// with container ids of all running docker containers
			err := dns.loadContainerInfo(cli, config.Hostname)

			if err == nil {
				// proceed after container was found
				done <- true
			}

			// wait 5 seconds to check again
			<-time.After(time.Second * 5)
		}
	}()

	// handle proceed after container was found or timeout after 30 seconds
	select {
	// the go routine to identify the kolumbus docker container will write on
	// the done channel to notify the main process, that the container was found
	// and it can proceed
	case <-done:
		log.Println("- container identified")
	// end identifying container after 30 seconds without success
	case <-time.After(time.Second * 30):
		errs <- errors.New("could not identify kolumbus container")
		return
	}
}

// loadContainerInfo will retrieve the container information of the kolumbus
// docker container
func (dns *Kolumbus) loadContainerInfo(cli *docker.Client, hostname string) error {

	// get a list of all running containers through the docker cli
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: false})
	if err != nil {
		return errors.Wrap(err, "could not fetch the list of docker containers")
	}

	for _, container := range containers {

		// check if the start of the container id matches the hostname. by default
		// docker creates the HOSTNAME environment variable equal to the start of
		// the docker container id
		if strings.HasPrefix(container.ID, hostname) {
			dns.ContainerInfo = &container
			return nil
		}
	}

	return errors.New("could not retrieve container information")
}
