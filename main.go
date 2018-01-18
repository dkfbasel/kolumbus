// kolumbus will watch docker for specific labels and create a dns server that
// enables envoyproxy to create a dynamic microservice mesh
package main

import "log"

func main() {

	kolumbus := FindABraveNewWorld()

	// start watching docker containers for services
	err := kolumbus.StartDockerWatch()
	if err != nil {
		log.Fatalf("could not init docker: %+v\n", err)
	}

	// start a dns service for envoyproxies to automatically
	// create a service mesh
	err = kolumbus.StartDNSServer()
	if err != nil {
		log.Fatalf("could not start server: %+v\n", err)
	}

}
