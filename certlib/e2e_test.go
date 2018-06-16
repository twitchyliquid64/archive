package certlib

import (
	"crypto/x509"
	"testing"
	"time"
)

func TestGenerateWorksBasicCert(t *testing.T) {
	ca := MakeBasicCA()
	fullCa, err := MakeCertKeyPair(ca)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	sub := MakeBasicCert()
	fullSub, err := MakeDerivedCertKeyPair(fullCa, sub)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = CheckCertSignatureIsDerived(fullCa.Cert, fullSub.Cert)
	if err != nil {
		t.Error(err)
	}
}

func TestGenerateWorksServerCert(t *testing.T) {
	ca := MakeBasicCA()
	fullCa, err := MakeCertKeyPair(ca)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	sub := MakeBasicCert()
	fullSub, err := MakeDerivedCertKeyPair(fullCa, sub)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = CheckCertSignatureIsDerived(fullCa.Cert, fullSub.Cert)
	if err != nil {
		t.Error(err)
	}

	rootCertPool := x509.NewCertPool()
	rootCertPool.AddCert(fullCa.Cert)

	opts := x509.VerifyOptions{
		Roots:         rootCertPool,
		CurrentTime:   time.Now(),
		Intermediates: x509.NewCertPool(),
	}

	chain, err := fullSub.Cert.Verify(opts)
	if err != nil {
		t.Error(err)
	}
	if len(chain) != 1 {
		t.Error("Expected 1 chain")
	}

	if len(chain[0]) != 2 {
		t.Error("expected 1 chain of 2 elements")
	}
}
