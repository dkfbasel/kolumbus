// kolumbus will watch docker for specific labels and create a dns server that
// enables envoyproxy to create a dynamic microservice mesh
package main

import "log"

func main() {

	kolumbus := FindABraveNewWorld()

	// initialize a channel of errors
	errorChan := make(chan error)

	// start watching docker containers for services
	kolumbus.StartDockerWatch(errorChan)
	log.Println("- docker container watch started")

	// start an envoyproxy and optionally open a port for
	// external communication
	kolumbus.StartEnvoyproxy(1494, 1498, errorChan)
	log.Println("- envoy proxy started")

	// start a data plane discovery service for envoyproxies to automatically
	// create a service mesh
	kolumbus.StartEnvoyDataPlaneServer(1492, errorChan)
	log.Println("- envoy discovery server started")

	// log any errors
	for err := range errorChan {
		log.Printf("%+v\n", err)
	}

	// proxy on server:
	// host and port to start the server on

	// proxy on remote machine:
	// host and port of remote address

}
