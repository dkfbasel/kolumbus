// kolumbus will watch docker for specific labels and create a dns server that
// enables envoyproxy to create a dynamic microservice mesh
package main

import "log"

func main() {

	kolumbus := NewKolumbus()

	err := kolumbus.StartDockerWatch()
	if err != nil {
		log.Fatalf("could not init docker: %+v\n", err)
	}

	err = kolumbus.StartServer()
	if err != nil {
		log.Fatalf("could not start server: %+v\n", err)
	}

}
