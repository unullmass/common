package tls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestVerifyCertBySha256(t *testing.T) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Organization:  []string{"Intel"},
			Country:       []string{"US"},
			Province:      []string{"CA"},
			Locality:      []string{"Folsom"},
			StreetAddress: []string{"1900 Prarie City Rd"},
			PostalCode:    []string{"95630"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey
	cert, _ := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)

	certFile, _ := ioutil.TempFile("", "cert.crt")
	defer os.Remove(certFile.Name())
	keyFile, _ := ioutil.TempFile("", "key.key")
	defer os.Remove(keyFile.Name())
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
	certFile.Close()
	pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyFile.Close()

	certDigest := sha256.Sum256(cert)

	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bar"))
	})
	go http.ListenAndServeTLS(":1337", certFile.Name(), keyFile.Name(), nil)

	tlsConfig := tls.Config{
		InsecureSkipVerify:    true,
		VerifyPeerCertificate: VerifyCertBySha256(certDigest),
	}
	transport := http.Transport{
		TLSClientConfig: &tlsConfig,
	}
	client := &http.Client{Transport: &transport}
	rsp, err := client.Get("https://localhost:1337/foo")
	if err != nil {
		t.Fatal(err)
	}
	rspBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Fatal(err)
	}
	rspString := string(rspBytes)
	if rspString != "bar" {
		t.Fail()
	}
}

func TestGenerateSelfSignCerts(t *testing.T) {

}

func TestGenRSAKeys(t *testing.T) {

}
