package main

import (
	"crypto/x509"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/twitchyliquid64/certlib"
)

var countryVar string
var organisationVar string
var organisationalUnitVar string
var expiryVar time.Duration
var makeCAVar bool
var oDerVar, oPemVar, oKeyPemVar, oKeyPkcsVar string
var commonNameVar, dnsNamesListVar string

func init() {
	flag.StringVar(&countryVar, "country", "U.S", "Country to display in the certificates metadata")
	flag.StringVar(&organisationVar, "org", "U.S", "Organisation to display in the certificates metadata")
	flag.StringVar(&organisationalUnitVar, "orgunit", "U.S", "Organisational unit to display in the certificates metadata")
	flag.DurationVar(&expiryVar, "expiry", time.Hour*time.Duration(24*30*6), "Time till expiry")
	flag.BoolVar(&makeCAVar, "ca", false, "Whether the generated cert is a trust root/intermediary")

	flag.StringVar(&oDerVar, "oDer", "", "Optional output filename for DER-formatted cert")
	flag.StringVar(&oPemVar, "oPEM", "", "Optional output filename for PEM-formatted cert")
	flag.StringVar(&oKeyPemVar, "oKeyPEM", "", "Optional output filename for PEM-formatted private key")
	flag.StringVar(&oKeyPkcsVar, "oKeyPkcs", "", "Optional output filename for PKCS-formatted private key")

	flag.StringVar(&commonNameVar, "commonName", "", "Certificate common name")
	flag.StringVar(&dnsNamesListVar, "dnsNames", "", "Comma-separated list of DNS names")
}

func main() {
	flag.Parse()

	if oDerVar == "" && oPemVar == "" && oKeyPemVar == "" && oKeyPkcsVar == "" {
		fmt.Println("Err: At least one output filename must be specified.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var cert *x509.Certificate
	if makeCAVar {
		cert = certlib.MakeBasicCA()
	} else {
		cert = certlib.MakeBasicServerCert()
	}
	certlib.SetDetails(cert, countryVar, organisationVar, organisationalUnitVar)

	cert.Subject.CommonName = commonNameVar
	cert.DNSNames = strings.Split(dnsNamesListVar, ",")

	pair, err := certlib.MakeCertKeyPair(cert)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	err = certlib.WritePrivateCertToFile(oDerVar, oPemVar, oKeyPemVar, oKeyPkcsVar, pair)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
