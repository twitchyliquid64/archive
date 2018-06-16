package certlib

import (
	"bytes"
	"crypto/x509"
	"encoding/binary"
	"math/big"
	unsecure_rand "math/rand"
	"testing"
	"time"
)

func TestGenerateRSAErrorsIfSmallKey(t *testing.T) {
	l, e := GenerateRSA(1024)
	if l != nil {
		t.Error("Got key, expected nil")
	}
	if e == nil {
		t.Error("Expected error")
	}
	if e != ErrInsecureKeyBitSize {
		t.Error("Expected InsecureKeyBitSizeErr")
	}
}

func TestMakeBasicCertWorks(t *testing.T) {
	testTime := time.Now()
	unsecure_rand.Seed(1)
	subjectKeyNum := uint64(unsecure_rand.Int63())
	var subjectKeyBytes = make([]byte, 16)
	binary.PutUvarint(subjectKeyBytes, subjectKeyNum)
	expectedSerialNumber := big.NewInt(int64(unsecure_rand.Int63()))
	unsecure_rand.Seed(1)

	cert := makeBasicCert(testTime)
	if cert.SerialNumber.String() != expectedSerialNumber.String() {
		t.Error("Expected serial number ", expectedSerialNumber.String(), ", got: ", cert.SerialNumber.String())
	}
	if bytes.Compare(cert.SubjectKeyId, subjectKeyBytes[:5]) != 0 {
		t.Error("Unexpected SubjectKeyId")
	}
	if cert.NotBefore != testTime {
		t.Error("NotBefore incorrect")
	}
	if cert.NotAfter != testTime.AddDate(0, 6, 0) {
		t.Error("Unexpected NotAfter")
	}
}

func TestMakeBasicCAWorks(t *testing.T) {
	unsecure_rand.Seed(1)
	cert := MakeBasicCert()
	makeBasicCA(cert)
	if cert.IsCA != true {
		t.Error("Expected IsCA to be set")
	}
	if cert.BasicConstraintsValid != true {
		t.Error("Expected BasicContraintsValid to be set")
	}
	if len(cert.ExtKeyUsage) != 2 || cert.ExtKeyUsage[0] != x509.ExtKeyUsageClientAuth || cert.ExtKeyUsage[1] != x509.ExtKeyUsageServerAuth {
		t.Error("Expected ExtKeyUsage to have ExtKeyUsageClientAuth & ExtKeyUsageServerAuth")
	}
	if cert.KeyUsage != (x509.KeyUsageCertSign) {
		t.Error("KeyUsage incorrect")
	}
}

func TestSetDetailsWorks(t *testing.T) {
	unsecure_rand.Seed(1)
	cert := MakeBasicCert()
	makeBasicCA(cert)
	SetDetails(cert, "Yololand", "Acme Incorporated", "Research & Development")
	if len(cert.Subject.Country) != 1 || cert.Subject.Country[0] != "Yololand" {
		t.Error("Country incorrect")
	}
	if len(cert.Subject.Organization) != 1 || cert.Subject.Organization[0] != "Acme Incorporated" {
		t.Error("Organization incorrect")
	}
	if len(cert.Subject.OrganizationalUnit) != 1 || cert.Subject.OrganizationalUnit[0] != "Research & Development" {
		t.Error("OrganizationalUnit incorrect")
	}
}
