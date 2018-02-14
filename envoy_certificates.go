package main

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// HandleEnvoyCertificateRequest will handle envoy requests for ssl certificates
func HandleEnvoyCertificateRequest(dns *Kolumbus, errs chan<- error) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// lock access to dns records
		dns.RLock()
		dns.RUnlock()

		// initialize a fingerprint to match client certificates
		// fingerprint := Fingerprint{}
		// fingerprint.Sha256 = "abdckdkediskdkfasdklfaskkdkjfkasdjf"

		response := CertificateResponse{}
		response.Certificates = []Fingerprint{}

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

// CertificateResponse is used for responses of a certificate request,
// it should provide a list of valid certificates
type CertificateResponse struct {
	Certificates []Fingerprint `json:"certificates"`
}

// Fingerprint is contains the sha256 hash of valid certificates
type Fingerprint struct {
	Sha256 string `json:"fingerprint_sha256,omitempty"`
}
