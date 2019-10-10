// kolumbus will watch docker for specific labels and create a dns server that
// enables envoyproxy to create a dynamic microservice mesh
package main

import (
	"log"
	"os"

	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

func main() {

	// initialize the configuration
	config := Config{}

	// get the container host name
	config.Hostname = os.Getenv("HOSTNAME")

	flag.IntVar(&config.DataPlanePort, "dataplane", 1492, "port to start the envoyproxy data plane discovery service on")
	flag.IntVar(&config.LocalProxyPort, "local-proxy", 1494, "port to start local proxy service of the internal envoyproxy instance on")
	flag.StringVar(&config.RemoteProxyMode, "remote-proxy-mode", MODE_NONE, "modus to use for remote proxy service (none/outbound/inbound)")
	flag.IntVar(&config.RemoteProxyPort, "remote-proxy-port", 1498, "(inbound) port to start the remote proxy service on the internal envoyproxy instance on")
	flag.StringVar(&config.RemoteProxyAddress, "remote-proxy-address", "", "(outbound) address of the remote proxy service to call")
	flag.Parse()

	// log the configuration
	log.Printf("- configuration: %+v\n", config)

	kolumbus := FindABraveNewWorld()

	// initialize a channel of errors
	errorChan := make(chan error)

	// log any errors in a separate go routine
	go func() {
		for err := range errorChan {
			log.Printf("%+v\n", err)
			if cause := errors.Cause(err); cause != nil {
				log.Printf("-- %+v\n", cause)
			}
		}
	}()

	kolumbus.IdentifyContainer(config, errorChan)
	log.Println("- start identifying container")

	// start watching docker containers for services
	kolumbus.StartDockerWatch(config, errorChan)
	log.Println("- docker container watch started")

	// start an envoyproxy and optionally open a port for
	// external communication
	kolumbus.StartEnvoyproxy(config, errorChan)
	log.Println("- envoy proxy started")

	// start a data plane discovery service for envoyproxies to automatically
	// create a service mesh
	kolumbus.StartEnvoyDataPlaneServer(config, errorChan)
	log.Println("- envoy discovery server started")

	// keep application running
	run := make(chan bool)
	<-run
}

// Config is used to define customizable configuration options
type Config struct {
	// define a port for envoy proxy data plane server that will be checked
	// by all envoyproxy instances for configuration information
	DataPlanePort int

	// define a port for a local envoyproxy instance that will be used
	// as local proxy for development
	LocalProxyPort int

	// define a port for a local envoyproxy instance that will be used for
	// remote connections to the cluster. connections are secured with tls
	// the mode (inbound/outbound) will define if envoyproxy will listen to
	// requests or send out requests to a remote cluster
	RemoteProxyMode    string // inbound or outbound
	RemoteProxyPort    int    // port to start the remote service on (inbound)
	RemoteProxyAddress string // address for a remote service to call (outbound)

	// hostname of the container. this is used to identify the corresponding
	// container in the docker container list
	Hostname string
}

// nolint
const MODE_NONE = "none"
const MODE_INBOUND = "inbound"
const MODE_OUTBOUND = "outbound"
