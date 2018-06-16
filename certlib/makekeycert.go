package certlib

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

// FullCert represents a valid certificate and private key combination.
type FullCert struct {
	Cert     *x509.Certificate
	DerBytes []byte
	Key      *rsa.PrivateKey
}

// MakeCertKeyPair generates a strong private key and signs the cert with it.
func MakeCertKeyPair(cert *x509.Certificate) (*FullCert, error) {
	var ret FullCert
	priv, err := GenerateRSA(2048)
	if err != nil {
		return nil, err
	}
	ret.Key = priv

	caBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}
	ret.DerBytes = caBytes
	caCert, err := x509.ParseCertificate(caBytes)
	if err != nil {
		return nil, err
	}
	ret.Cert = caCert
	return &ret, nil
}

// MakeDerivedCertKeyPair generates a strong private key and signs the cert with it.
// It then signs the cert with the CA cert provided, such that the generated cert
// can be proven to be associated with the CA cert.
func MakeDerivedCertKeyPair(ca *FullCert, cert *x509.Certificate) (*FullCert, error) {
	var ret FullCert
	priv, err := GenerateRSA(2048)
	if err != nil {
		return nil, err
	}
	ret.Key = priv

	caBytes, err := x509.CreateCertificate(rand.Reader, cert, ca.Cert, &priv.PublicKey, ca.Key)
	if err != nil {
		return nil, err
	}
	ret.DerBytes = caBytes
	caCert, err := x509.ParseCertificate(caBytes)
	if err != nil {
		return nil, err
	}
	ret.Cert = caCert
	return &ret, nil
}
