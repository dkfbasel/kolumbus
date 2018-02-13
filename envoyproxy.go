package main

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// StartEnvoyproxy will initialize an envoy proxy that can be used to connect
// services locally or with a remote service (local services will be prioritized)
func (dns *Kolumbus) StartEnvoyproxy(errs chan<- error) {

	// command to start envoy
	// /usr/local/bin/envoy -c /app/envoy/config.json --v2-config-only -l error

	go func() {
		bin := "/usr/local/bin/envoy"
		args := []string{"-c", "/app/envoy/config/proxy.json", "--v2-config-only", "-l error"}

		cmd := exec.Command(bin, args...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			errs <- errors.Wrap(err, "could not start envoy proxy")
		}
	}()

}
