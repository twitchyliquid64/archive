package certlib

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

// LoadPrivateCert returns a certificate and private key, decoded from bytesCert (DER / ASN format) and keyBytes (ASN1).
func LoadPrivateCert(bytesCert []byte, keyBytes []byte) (*FullCert, error) {
	cert, err := x509.ParseCertificate(bytesCert)
	if err != nil {
		return nil, err
	}
	priv, err := x509.ParsePKCS1PrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}
	return &FullCert{
		Cert:     cert,
		Key:      priv,
		DerBytes: bytesCert,
	}, nil
}

// LoadPrivateCertPEM returns a certificate and private key, decoded from bytesCert (PEM) and keyBytes (PEM).
func LoadPrivateCertPEM(bytesCert []byte, keyBytes []byte) (*FullCert, error) {
	certDERBlock, _ := pem.Decode(bytesCert)
	if certDERBlock == nil {
		return nil, errors.New("No certificate data read from PEM")
	}
	cert, err := x509.ParseCertificate(certDERBlock.Bytes)
	if err != nil {
		return nil, err
	}

	keyBlock, _ := pem.Decode(keyBytes)
	if keyBlock == nil {
		return nil, errors.New("No key data read from PEM")
	}
	priv, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return &FullCert{
		Cert:     cert,
		Key:      priv,
		DerBytes: certDERBlock.Bytes,
	}, nil
}

// LoadPrivateCertFromFilePEM returns a cert & PK after loading both those components from the files at the specified paths.
// certPath should point to a PEM encoded certificate, and keyPath should point to a PEM encoded private key.
func LoadPrivateCertFromFilePEM(certPath, keyPath string) (*FullCert, error) {
	certBytes, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return LoadPrivateCertPEM(certBytes, keyBytes)
}

// LoadPrivateCertFromFile returns a cert & PK after loading both those components from the files at the specified paths.
// certPath should point to a DER/ASN1 encoded certificate, and keyPath should point to a ASN1 encoded private key.
func LoadPrivateCertFromFile(certPath, keyPath string) (*FullCert, error) {
	certBytes, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return LoadPrivateCert(certBytes, keyBytes)
}

// WritePrivateCertToFile writes a full certificate of yours to a set of files, so it can be later loaded.
// passing "" to any of the paths omits that file's generation.
//
func WritePrivateCertToFile(derFile, certPEMFile, keyPEMFile, keyPKCSFile string, cert *FullCert) error {
	if derFile != "" {
		certCerFile, err := os.Create(derFile)
		if err != nil {
			return err
		}
		certCerFile.Write(cert.DerBytes)
		certCerFile.Close()
	}

	if certPEMFile != "" {
		certFile, err := os.Create(certPEMFile)
		if err != nil {
			return err
		}
		pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: cert.DerBytes})
		certFile.Close()
	}

	if keyPEMFile != "" {
		keyFile, err := os.Create(keyPEMFile)
		if err != nil {
			return err
		}
		pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(cert.Key)})
		keyFile.Close()
	}

	if keyPKCSFile != "" {
		keyFile, err := os.Create(keyPKCSFile)
		if err != nil {
			return err
		}
		privBytes := x509.MarshalPKCS1PrivateKey(cert.Key)
		keyFile.Write(privBytes)
		keyFile.Close()
	}

	return nil
}
