package certlib

import "crypto/x509"

// CheckCertSignatureIsDerived returns error == nil if suspectCert is derived/signed by caCert.
func CheckCertSignatureIsDerived(caCert *x509.Certificate, suspectCert *x509.Certificate) error {
	return suspectCert.CheckSignatureFrom(caCert)
}
