package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

const certificateDir = "/app/envoy/certificates"

// StartEnvoyproxy will initialize an envoy proxy that can be used to connect
// services locally or with a remote service (local services will be prioritized)
func (dns *Kolumbus) StartEnvoyproxy(config Config, errs chan<- error) {

	// command to start envoy
	// /usr/local/bin/envoy -c /app/envoy/config.json --v2-config-only --log-level error

	go func() {

		// handle remote proxy modes differently
		// note: proxy will only be enabled if the configuration is correct
		switch config.RemoteProxyMode {
		case MODE_OUTBOUND:
			// outbound remote proxies need a client certificate that must be
			// provided from the host volume as it needs to be signed by the
			// corresponding server certificate on the remote host

			// a remote proxy address must be given
			if config.RemoteProxyAddress == "" {
				config.RemoteProxyMode = MODE_NONE
				errs <- errors.New("A remote proxy address and port must be specified for outbound proxy mode")
			}

			// a client certificate must be provided
			_, err := os.Stat(filepath.Join(certificateDir, "client.crt"))
			if os.IsNotExist(err) {
				config.RemoteProxyMode = MODE_NONE
				errs <- errors.New("A client certificate (client.crt) must be specified for connections to a remote cluster")
			}

			// a client key must be provided
			_, err = os.Stat(filepath.Join(certificateDir, "client.key"))
			if os.IsNotExist(err) {
				config.RemoteProxyMode = MODE_NONE
				errs <- errors.New("A client key (client.key) must be specified for connections to a remote cluster")
			}

		case MODE_INBOUND:
			// inbound remote proxies need certificates for ssl connections
			// we will automatically create respective certificats if not yet exisiting

			// generate a server certificate if necessary
			_, err := os.Stat(filepath.Join(certificateDir, "server.crt"))
			if os.IsNotExist(err) {

				log.Println("- generating server certificate")

				// generate a server certificate and key
				cmd := exec.Command("openssl",
					"req", "-x509", "-newkey", "rsa:4096",
					"-keyout", "server.key",
					"-out", "server.crt",
					"-nodes", "-days", "365",
					"-subj", "/CN=kolumbus.remote",
				)
				cmd.Dir = certificateDir
				err = cmd.Run()
				if err != nil {
					errs <- errors.Wrap(err, "could not generate server certificate")
					config.RemoteProxyMode = MODE_NONE
					break
				}
			}

			// generate a client certificate if necessary
			_, err = os.Stat(filepath.Join(certificateDir, "client.crt"))
			if os.IsNotExist(err) {

				log.Println("- generating client certificate")

				// TODO: archive any existing keys

				// generate a client certificate and certificate signign request
				cmd := exec.Command("openssl",
					"req", "-newkey", "rsa:4096",
					"-keyout", "client.key",
					"-out", "client.csr",
					"-nodes", "-days", "365",
					"-subj", "/CN=kolumbus.client",
				)
				cmd.Dir = certificateDir
				err = cmd.Run()
				if err != nil {
					errs <- errors.Wrap(err, "could not generate client certificate signign request")
					config.RemoteProxyMode = MODE_NONE
					break
				}

				// sign the client certificate signing requesg with the server
				// certificate to create a client certificate
				cmd2 := exec.Command("openssl",
					"x509", "-req",
					"-in", "client.csr",
					"-CA", "server.crt",
					"-CAkey", "server.key",
					"-out", "client.crt",
					"-set_serial", "01",
					"-days", "365")
				cmd2.Dir = certificateDir
				err = cmd2.Run()
				if err != nil {
					errs <- errors.Wrap(err, "could not sign client certificate")
					config.RemoteProxyMode = MODE_NONE
				}

				// delete the certificate signign request
				_ = os.Remove(filepath.Join(certificateDir, "client.csr"))
			}

		}

		// parse the envoy template configuration configuration
		tmpl, err := template.ParseFiles("/app/envoy/config/envoyproxy.config.tmpl")
		if err != nil {
			errs <- errors.Wrap(err, "could not parse envoyproxy configuration template")
			return
		}

		// create a new temporary files for the envoyproxy configuration
		file, err := os.Create("/tmp/envoyproxy.config.json")
		if err != nil {
			errs <- errors.Wrap(err, "could not open envoyproxy configuration file for writing")
			return
		}

		// compile the template with the configuration
		err = tmpl.Execute(file, config)
		file.Close() //nolint:errcheck

		if err != nil {
			errs <- errors.Wrap(err, "could not write envoyproxy configuration file")
			return
		}

		// start envoyproxy with the given configuration

		bin := "/usr/local/bin/envoy"
		args := []string{"-c", "/tmp/envoyproxy.config.json", "--v2-config-only", "--log-level", "error"}

		cmd := exec.Command(bin, args...)

		// redirect output of envoyproxy to stdout and stderr
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// run envoyproxy and wait for it to finish
		err = cmd.Run()
		if err != nil {
			errs <- errors.Wrap(err, "could not start envoy proxy")
		}
	}()

}
