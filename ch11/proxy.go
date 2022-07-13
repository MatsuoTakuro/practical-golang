package ch11

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func proxy() {
	// ssl()
	InsecureSkipVerify()
}

func ssl() {

	// Generating a certificate
	// ref. https://www.section.io/engineering-education/how-to-get-ssl-https-for-localhost/

	cert, err := os.ReadFile("ch11/.key/localhost.crt")
	if err != nil {
		panic(err)
	}

	// X.509 means the format of public key certificates defined by ITU https://en.wikipedia.org/wiki/X.509#:~:text=An%20X.,authority%20or%20is%20self%2Dsigned.
	certPool := x509.NewCertPool() // CertPool is a set of certificates.
	certPool.AppendCertsFromPEM(cert)
	cfg := &tls.Config{
		RootCAs: certPool,
	}
	cfg.BuildNameToCertificate()

	hc := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: cfg,
			Proxy:           http.ProxyFromEnvironment,
		},
	}

	// TODO: panic: Get "https://www.oreilly.co.jp/index.shtml": x509: certificate signed by unknown authority
	resp, err := hc.Get("https://www.oreilly.co.jp/index.shtml")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func InsecureSkipVerify() {
	hc := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: http.ProxyFromEnvironment,
		},
	}
	resp, err := hc.Get("https://www.oreilly.co.jp/index.shtml")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
