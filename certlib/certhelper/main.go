package main

import (
	"bufio"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/twitchyliquid64/certlib"
)

var modeVar string

const (
	caCertFilePath    = "ca.pem"
	caDepCertFilePath = "ca.crt"
	caKeyFilePath     = "ca.priv"
)

func init() {
	flag.StringVar(&modeVar, "mode", "", "What certhelper should do: makeCA/showCA/forgeCert/showCert/combineCerts")
}

func main() {
	flag.Parse()

	if modeVar == "makeCA" {
		makeCA()
	} else if modeVar == "forgeCert" {
		forgeCert()
	} else if modeVar == "showCA" {
		showCA()
	} else if modeVar == "showCert" {
		showCert()
	} else if modeVar == "combineCerts" {
		name := stringPrompt("Please enter the cert name, excluding file extension: ", "", true)
		outname := stringPrompt("Please enter the output cert name, excluding file extension: ", "", true)

		if !fileExists(name + ".certPEM") {
			fmt.Println("Err: could not find file " + name + ".certPEM")
			os.Exit(1)
		}
		if !fileExists(caCertFilePath) {
			fmt.Println("Err: could not find " + caCertFilePath)
			os.Exit(1)
		}
		out, err := exec.Command("sh", "-c", "cat "+name+".certPEM "+caCertFilePath+" > "+outname+".certPEM").Output()
		fmt.Println(string(out))
		if err != nil {
			fmt.Println("Error running command: ", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Error: a Valid mode must be specified.")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func booleanPrompt(promptText, yesResponse string) bool {
	fmt.Print(promptText)
	var input string
	fmt.Scanln(&input)
	return input == yesResponse
}

func stringPrompt(promptText, defaultValue string, needValue bool) string {
	var input string
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(promptText)
	input, _ = reader.ReadString('\n')
	for strings.TrimSpace(input) == "" && needValue {
		fmt.Print(promptText)
		input, _ = reader.ReadString('\n')
	}
	if input == "" {
		input = defaultValue
	}
	return strings.TrimSpace(input)
}

func showCert() {
	name := stringPrompt("Please enter the cert name, excluding file extension: ", "", true)
	if !fileExists(name + ".certPEM") {
		fmt.Println("Err: could not find file " + name + ".certPEM")
		os.Exit(1)
	}
	certBytes, err := ioutil.ReadFile(name + ".certPEM")
	if err != nil {
		fmt.Println("Err: ", err)
		os.Exit(1)
	}
	certDERBlock, _ := pem.Decode(certBytes)
	if certDERBlock == nil {
		fmt.Println("Error: No cert data read from PEM")
		os.Exit(1)
	}
	cert, err := x509.ParseCertificate(certDERBlock.Bytes)
	if err != nil {
		fmt.Println("Error parsing cert: ", err)
		os.Exit(1)
	}
	printCertDetails(cert)
}

func forgeCert() {
	ca := readCaFromDisk()
	fmt.Println("CA:")
	printCertDetails(ca.Cert)

	cert := certlib.MakeBasicServerCert()
	promptForCertInfo(cert, ca.Cert)

	cert.Issuer = ca.Cert.Subject

	fullSubjectCert, err := certlib.MakeDerivedCertKeyPair(ca, cert)
	if err != nil {
		fmt.Println("Err: ", err)
		os.Exit(1)
	}

	fmt.Println("Writing cert to " + cert.Subject.CommonName + ".certPEM")
	fmt.Println("Writing key to " + cert.Subject.CommonName + ".keyPEM")
	err = certlib.WritePrivateCertToFile("", cert.Subject.CommonName+".certPEM", cert.Subject.CommonName+".keyPEM", "", fullSubjectCert)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func showCA() {
	ca := readCaFromDisk()
	fmt.Println("CA:")
	printCertDetails(ca.Cert)
}

func makeCA() {
	if fileExists(caCertFilePath) || fileExists(caKeyFilePath) {
		if !booleanPrompt("A CA exists (ca.pem, ca.priv), overwrite? [y/N]: ", "y") {
			fmt.Println("Abort.")
			os.Exit(1)
		}
	}
	cert := certlib.MakeBasicCA()

	promptForCertInfo(cert, nil)

	pair, err := certlib.MakeCertKeyPair(cert)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	key := stringPrompt("Enter a strong password to encrypt the key on disk: ", "", true)

	fmt.Println("Writing cert to " + caCertFilePath + " & " + caDepCertFilePath)
	fmt.Println("Writing encrypted key to " + caKeyFilePath)
	certlib.WritePrivateCertToFile(caDepCertFilePath, caCertFilePath, "", "", pair)
	ciphertext := encrypt(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pair.Key)}), key)
	ioutil.WriteFile(caKeyFilePath, ciphertext, 0600)
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func readCaFromDisk() *certlib.FullCert {
	if !fileExists(caCertFilePath) || !fileExists(caKeyFilePath) {
		fmt.Println("Error: Could not stat ca.pem & ca.priv")
		os.Exit(1)
	}

	key := stringPrompt("CA Key password: ", "", true)
	ciphertext, err := ioutil.ReadFile(caKeyFilePath)
	if err != nil {
		fmt.Println("Error reading ca.priv: ", err)
		os.Exit(1)
	}
	keyBlock, _ := pem.Decode(decrypt(ciphertext, key))
	if keyBlock == nil {
		fmt.Println("Error: No key data read from PEM")
		fmt.Println("Is the key correct?")
		os.Exit(1)
	}
	priv, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		fmt.Println("Error parsing PKCS1: ", err)
		os.Exit(1)
	}

	pemBytes, err := ioutil.ReadFile(caCertFilePath)
	if err != nil {
		fmt.Println("Error reading ca.pem: ", err)
		os.Exit(1)
	}
	certDERBlock, _ := pem.Decode(pemBytes)
	if certDERBlock == nil {
		fmt.Println("Error: No cert data read from PEM")
		os.Exit(1)
	}
	cert, err := x509.ParseCertificate(certDERBlock.Bytes)
	if err != nil {
		fmt.Println("Error parsing cert: ", err)
		os.Exit(1)
	}
	caPair := &certlib.FullCert{
		Cert:     cert,
		Key:      priv,
		DerBytes: certDERBlock.Bytes,
	}
	return caPair
}

func promptForCertInfo(cert, template *x509.Certificate) {
	var country string
	if template != nil && len(template.Subject.Country) > 0 {
		country = stringPrompt("What country should be written? ["+template.Subject.Country[0]+"]: ", template.Subject.Country[0], false)
	} else {
		country = stringPrompt("What country should be written? [U.S]: ", "U.S", false)
	}
	var org string
	if template != nil && len(template.Subject.Organization) > 0 {
		org = stringPrompt("What organisation should be written? ["+template.Subject.Organization[0]+"]: ", template.Subject.Organization[0], false)
	} else {
		org = stringPrompt("What organisation should be written? [Acme Co.]: ", "Acme Co.", false)
	}
	var orgU string
	if template != nil && len(template.Subject.OrganizationalUnit) > 0 {
		orgU = stringPrompt("What org unit should be written? ["+template.Subject.OrganizationalUnit[0]+"]: ", template.Subject.OrganizationalUnit[0], false)
	} else {
		orgU = stringPrompt("What org unit should be written? [Acme Certs]: ", "Acme Certs.", false)
	}
	commonName := stringPrompt("Whats the common name/domain [example.com]: ", "example.com", false)
	dns := stringPrompt("Enter addition domains in CSV form []: ", "", false)

	certlib.SetDetails(cert, country, org, orgU)
	cert.Subject.CommonName = commonName
	cert.DNSNames = strings.Split(dns, ",")
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
