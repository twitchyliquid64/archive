package tls

import (
	"crypto/x509"
	"errors"
	"log"
	"time"
)

// VerifyPeerCertificateFunc is a shorthand for the type expected by tls.Config.VerifyPeerCertificate.
type VerifyPeerCertificateFunc func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error

// CheckCertSignature returns a function which will check the signature and expiry of the
// certificates it is presented with. A <nil> will only be returned if a certificate presented
// was signed by the CA certificate, and is not expired.
func CheckCertSignature(caCert *x509.Certificate) VerifyPeerCertificateFunc {
	return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		for _, c := range rawCerts {
			parsedCert, err := x509.ParseCertificate(c)
			if err != nil {
				return err
			}
			certErr := parsedCert.CheckSignatureFrom(caCert)
			if parsedCert.NotAfter.Before(time.Now()) || parsedCert.NotBefore.After(time.Now()) {
				certErr = errors.New("Certificate expired or used too soon")
			}
			log.Printf("Remote presented certificate %d with time bounds (%v-%v). Verification error for certificate: %+v",
				parsedCert.SerialNumber, parsedCert.NotBefore, parsedCert.NotAfter, certErr)
			return certErr
		}
		return errors.New("Expected certificate which would pass, none presented")
	}
}
