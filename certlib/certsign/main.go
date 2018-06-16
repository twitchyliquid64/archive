package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/twitchyliquid64/certlib"
)

var noPrompt bool
var authorityPEMInputVar bool
var authorityCertVar, authorityKeyVar string
var subjectPEMInputVar bool
var subjectCertVar string
var oDerVar, oPemVar, oKeyPemVar, oKeyPkcsVar string

func init() {
	flag.BoolVar(&noPrompt, "y", false, "If set, will not prompt for confirmation")
	flag.BoolVar(&authorityPEMInputVar, "authorityIsPEM", false, "Set this if the authorities key/cert is in PEM format")
	flag.BoolVar(&subjectPEMInputVar, "subjectIsPEM", false, "Set this if the cert you are signing is in PEM format")

	flag.StringVar(&authorityCertVar, "authorityCert", "", "Path to authority's certificate")
	flag.StringVar(&authorityKeyVar, "authorityKey", "", "Path to authority's private key")

	flag.StringVar(&subjectCertVar, "cert", "", "Path to cert you wish to sign")

	flag.StringVar(&oDerVar, "oDer", "", "Optional output filename for DER-formatted cert")
	flag.StringVar(&oPemVar, "oPEM", "", "Optional output filename for PEM-formatted cert")
	flag.StringVar(&oKeyPemVar, "oKeyPEM", "", "Optional output filename for PEM-formatted private key")
	flag.StringVar(&oKeyPkcsVar, "oKeyPkcs", "", "Optional output filename for PKCS-formatted private key")
}

func checkFlags() {
	if authorityCertVar == "" || authorityKeyVar == "" {
		fmt.Println("Err: Both the private key and the certificate of the signing authority must be specified.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if subjectCertVar == "" || (oDerVar == "" && oPemVar == "") || (oKeyPemVar == "" && oKeyPkcsVar == "") {
		fmt.Println("Err: The certificate you wish to sign (-cert) must be specified, along with at least one output filename for both key/cert (Or both to write it out in both formats).")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	checkFlags()

	var err error
	var dataBytes []byte

	dataBytes, err = ioutil.ReadFile(subjectCertVar)
	if err != nil {
		fmt.Println("Err loading cert: ", err)
		os.Exit(1)
	}

	if subjectPEMInputVar {
		certDERBlock, _ := pem.Decode(dataBytes)
		if certDERBlock == nil {
			fmt.Println("Err: No certificate data read from cert PEM")
			os.Exit(1)
		}
		dataBytes = certDERBlock.Bytes
	}

	cert, err := x509.ParseCertificate(dataBytes)
	if err != nil {
		fmt.Println("Err parsing cert: ", err)
		os.Exit(1)
	}

	var CAcert *certlib.FullCert
	if authorityPEMInputVar {
		CAcert, err = certlib.LoadPrivateCertFromFilePEM(authorityCertVar, authorityKeyVar)
	} else {
		CAcert, err = certlib.LoadPrivateCertFromFile(authorityCertVar, authorityKeyVar)
	}
	if err != nil {
		fmt.Println("Err when loading authority cert/key: ", err)
		os.Exit(1)
	}

	if !noPrompt {
		fmt.Println("You are signing the certificate:")
		printCertDetails(cert)
		fmt.Println()
		fmt.Println("Using the certificate:")
		printCertDetails(CAcert.Cert)
		fmt.Println()
		fmt.Print("Are you sure you want to proceed [y/N]: ")
		var input string
		fmt.Scanln(&input)
		if len(input) == 0 || string(strings.ToUpper(input)[0]) != "Y" {
			fmt.Println("Aborting.")
			os.Exit(0)
		}
	}

	cert.Issuer = CAcert.Cert.Subject

	fullSubjectCert, err := certlib.MakeDerivedCertKeyPair(CAcert, cert)
	if err != nil {
		fmt.Println("Err: ", err)
		os.Exit(1)
	}

	err = certlib.WritePrivateCertToFile(oDerVar, oPemVar, oKeyPemVar, oKeyPkcsVar, fullSubjectCert)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func printCertDetails(cert *x509.Certificate) {
	fmt.Printf("  Serial %s, expiring at %s\n", cert.SerialNumber.String(), cert.NotAfter.String())
	if cert.Subject.CommonName != "" {
		fmt.Printf("  CN: %s\n", cert.Subject.CommonName)
	}
	if len(cert.Subject.Country) != 0 {
		fmt.Printf("  Country: %+v\n", cert.Subject.Country)
	}
	if len(cert.Subject.Organization) != 0 {
		fmt.Printf("  Organization: %+v\n", cert.Subject.Organization)
	}
	if len(cert.Subject.OrganizationalUnit) != 0 {
		fmt.Printf("  OrganizationalUnit: %+v\n", cert.Subject.OrganizationalUnit)
	}
	if cert.Subject.SerialNumber != "" {
		fmt.Printf("  Subject SN: %s\n", cert.Subject.SerialNumber)
	}
	if len(cert.Subject.ExtraNames) != 0 {
		fmt.Printf("  Extra Names: %+v\n", cert.Subject.ExtraNames)
	}
	if len(cert.DNSNames) != 0 {
		fmt.Printf("  DNS Names: %+v\n", cert.DNSNames)
	}
	if len(cert.IPAddresses) != 0 {
		fmt.Printf("  IP Addresses: %+v\n", cert.IPAddresses)
	}
	if len(cert.EmailAddresses) != 0 {
		fmt.Printf("  Email Addresses: %+v\n", cert.EmailAddresses)
	}
	fmt.Printf("  CA: %t\n", cert.IsCA)
}
