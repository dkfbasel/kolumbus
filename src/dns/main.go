// dockerdns will watch docker for specific labels and create a dns server that
// enables envoyproxy to create a dynamic microservice mesh
package main

import "log"

func main() {

	dns := NewDNS()

	err := dns.InitDocker()
	if err != nil {
		log.Fatalf("could not init docker: %+v\n", err)
	}

	err = dns.StartServer()
	if err != nil {
		log.Fatalf("could not start server: %+v\n", err)
	}

}
