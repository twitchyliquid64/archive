package certlib

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"errors"
	"math/big"
	unsecure_rand "math/rand"
	"time"
)

// ErrInsecureKeyBitSize is returned if a generate method is called with too few bits.
var ErrInsecureKeyBitSize = errors.New("too few bits when generating key")

// GenerateRSA returns a RSA private key with the given key length.
func GenerateRSA(bitSize int) (*rsa.PrivateKey, error) {
	if bitSize <= 1024 {
		return nil, ErrInsecureKeyBitSize
	}

	return rsa.GenerateKey(rand.Reader, bitSize)
}

// MakeBasicCert returns a basic x509 certificate with minimum fields set to sensible defaults. It is
// expected that users will further modify the certificate before it is used.
func MakeBasicCert() *x509.Certificate {
	//Use a different random number generator so we dont leak any state of crypto/rand
	//Who cares if the serial number is predictable, we know when the cert generated anyway
	//through NotBefore.
	unsecure_rand.Seed(time.Now().Unix())
	return makeBasicCert(time.Now())
}

func makeBasicCert(now time.Time) *x509.Certificate {
	//Make a subjectKeyId. There are no security requirements for this field, but the
	//more statistically distributed it is the better it can be used.
	subjectKeyNum := uint64(unsecure_rand.Int63())
	var subjectKeyBytes = make([]byte, 16)
	binary.PutUvarint(subjectKeyBytes, subjectKeyNum)

	return &x509.Certificate{
		SerialNumber: big.NewInt(int64(unsecure_rand.Int63())),
		Subject: pkix.Name{
			Country:            []string{"U.S"},
			Organization:       []string{"Acme Co."},
			OrganizationalUnit: []string{"Acme Co." + "U"},
		},
		NotBefore:    now,
		NotAfter:     now.AddDate(0, 6, 0), //6 month expiry
		SubjectKeyId: subjectKeyBytes[:5],
	}
}

// MakeBasicServerCert returns a basic x509 certificate with minimum fields necessary to act as a TLS / other server.
// It is expected the caller will set their own fields, like country / subject using SetDetails().
func MakeBasicServerCert() *x509.Certificate {
	cert := MakeBasicCert()
	makeBasicServerCert(cert)
	return cert
}

func makeBasicServerCert(cert *x509.Certificate) {
	cert.IsCA = false
	cert.BasicConstraintsValid = true
	cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	cert.KeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
}

// MakeBasicCA returns a basic x509 certificate with minimum fields necessary to act as a CA certificate.
// It is expected the caller will set their own fields, like country / subject.
func MakeBasicCA() *x509.Certificate {
	cert := MakeBasicCert()
	makeBasicCA(cert)
	return cert
}

func makeBasicCA(cert *x509.Certificate) {
	cert.IsCA = true
	cert.BasicConstraintsValid = true
	cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	cert.KeyUsage |= x509.KeyUsageCertSign
}

// SetDetails sets the human-readable details of a certificate. Use this after generating a certificate with an above method.
func SetDetails(cert *x509.Certificate, country, organisation, organisationUnit string) {
	cert.Subject = pkix.Name{
		Country:            []string{country},
		Organization:       []string{organisation},
		OrganizationalUnit: []string{organisationUnit},
	}
}
