# certlib
Helper library for working with certificates.

[![GoDoc](https://godoc.org/github.com/twitchyliquid64/certlib?status.svg)](http://godoc.org/github.com/twitchyliquid64/certlib) [![Build Status](https://travis-ci.org/twitchyliquid64/certlib.svg?branch=master)](https://travis-ci.org/twitchyliquid64/certlib) [![Go Report Card](https://goreportcard.com/badge/github.com/twitchyliquid64/certlib)](https://goreportcard.com/report/github.com/twitchyliquid64/certlib)


### `./certhelper`

The easy way is to invoke `certhelper`. This nifty cmdline utility will do everything you want in an interactive fashion, with the added benefit that your CA key
never touches disk unencrypted.

#### Make a CA

```
./certhelper --mode makeCA
... A few questions and answers ...
Writing cert to ca.pem & ca.crt
Writing encrypted key to ca.priv
```

#### Make a SSL certificate signed by your CA

```
./certhelper --mode forgeCert
CA Key password: ******
... A few questions and answers ...
Writing cert to thing.certPEM
Writing key to thing.keyPEM
```

#### Make the certificate chain for your SSL server

```shell
./certhelper --mode combineCerts
Please enter the cert name, excluding file extension: <certificate common name>
Please enter the output cert name, excluding file extension: combined
// Will produce combined.certPEM, which you can give to your SSL/HTTPS server
// (in Go, using tls.Config.Certificates). Dont throw away <SSL cert>.keyPEM either,
// you will need that as well.
```

### certgen & certsign

_NOTE: Run these commands with no arguments to see help/USAGE information._

Usage of these two utilities over `certhelper` is not recommended because of the number of complicated command line parameters.

Make a CA: `certgen -ca -org Example -orgunit "Example Key Authority" -ca -oDer CA.der -oKeyPkcs CA.key -oPEM CA.certPEM -oKeyPEM CA.keyPEM -commonName ExampleCA`

Create SSL certificate: `certgen -org Example -orgunit "Example Key Authority" -oPEM example.certPEM -oKeyPEM example.keyPEM -commonName example.com -dnsNames example.com`

Sign the SSL cert with your CA: `certsign -subjectIsPEM -authorityCert CA.der -authorityKey CA.key -cert example.certPEM -oPEM example.certPEM -oKeyPEM example.keyPEM`

NOTE: `tls.LoadX509KeyPair` requires paths to PEM encoded files - hence why I only bothered outputting the PEM files in the signing stage.

Another gotcha: if you intend to serve this SSL cert via `tls.Config`, be aware that you need to serve BOTH the CA cert and the SSL cert. This is easier than it sounds: `cat mycert.certPEM CA.certPEM > combined.certPEM`. Simply load this combined certificate into `config.Certificates` and you're ready to go.

### Installing a CA cert into linux

`certutil -d sql:$HOME/.pki/nssdb -A -t TC -n "<give your CA a name>" -i ca.crt `

You're welcome.

### TODO

- [ ] Moar tests
- [ ] Examples on the README
- [ ] Support generating other types of keys (ECC)
- [ ] Generic encryption
- [ ] Convert to TLS certificate type
- [ ] Annotate with allowed DNS names
- [ ] Annotate with custom ExtKeyUsage
- [ ] More helper methods for TLS
- [ ] Refactor some of the code in the two utilities into methods, in this package
