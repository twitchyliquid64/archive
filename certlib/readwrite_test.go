package certlib

import (
	"bytes"
	"os"
	"testing"
)

func TestStoreLoadDerWithPCKS(t *testing.T) {
	ca := MakeBasicCA()
	fullCa, err := MakeCertKeyPair(ca)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = WritePrivateCertToFile("testDer", "", "", "test.key", fullCa)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fullReadBack, err := LoadPrivateCertFromFile("testDer", "test.key")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if fullReadBack.Key.D.Cmp(fullCa.Key.D) != 0 || fullReadBack.Key.N.Cmp(fullCa.Key.N) != 0 || fullReadBack.Key.E != fullReadBack.Key.E {
		t.Error("Key mismatch")
	}
	if fullReadBack.Cert.SerialNumber.Cmp(fullCa.Cert.SerialNumber) != 0 {
		t.Error("Incorrect cert")
	}
	if bytes.Compare(fullReadBack.Cert.Signature, fullCa.Cert.Signature) != 0 {
		t.Error("Incorrect signature")
	}

	os.Remove("testDer")
	os.Remove("test.key")
}

func TestStoreLoadPEM(t *testing.T) {
	ca := MakeBasicCA()
	fullCa, err := MakeCertKeyPair(ca)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = WritePrivateCertToFile("", "cert.pem", "key.pem", "", fullCa)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fullReadBack, err := LoadPrivateCertFromFilePEM("cert.pem", "key.pem")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if fullReadBack.Key.D.Cmp(fullCa.Key.D) != 0 || fullReadBack.Key.N.Cmp(fullCa.Key.N) != 0 || fullReadBack.Key.E != fullReadBack.Key.E {
		t.Error("Key mismatch")
	}
	if fullReadBack.Cert.SerialNumber.Cmp(fullCa.Cert.SerialNumber) != 0 {
		t.Error("Incorrect cert")
	}
	if bytes.Compare(fullReadBack.Cert.Signature, fullCa.Cert.Signature) != 0 {
		t.Error("Incorrect signature")
	}

	os.Remove("cert.pem")
	os.Remove("key.pem")
}
