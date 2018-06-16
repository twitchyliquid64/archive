package tls

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

// ConfigureTLS generates and returns a TLS configuration based on the given parameters.
func ConfigureTLS(certPemPath, keyPemPath, caCertPath string) (*tls.Config, error) {
	var caCertParsed *x509.Certificate
	pemBytes, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}
	certDERBlock, _ := pem.Decode(pemBytes)
	if certDERBlock == nil {
		return nil, errors.New("No certificate data read from PEM")
	}
	caCertParsed, err = x509.ParseCertificate(certDERBlock.Bytes)
	if err != nil {
		return nil, err
	}

	c := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		VerifyPeerCertificate: CheckCertSignature(caCertParsed),
		InsecureSkipVerify:    true,
		ClientAuth:            tls.RequestClientCert,
	}

	mainCert, err := tls.LoadX509KeyPair(certPemPath, keyPemPath)
	if err != nil {
		return nil, err
	}
	c.Certificates = []tls.Certificate{mainCert}

	return c, nil
}
