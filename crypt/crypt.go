package crypt

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"time"
)

// SignedData contains the original byte stream that was signed along with
type SignedData struct {
	Data      []byte `json:"data"` //json formatted VMtrust report.
	Alg       string `json:"hash_alg"`
	Cert      string `json:"cert"` //pem formatted certificate
	Signature []byte `json:"signature"`
}

// GetHexRandomString return a random string of 'length'
func GetHexRandomString(length int) (string, error) {

	bytes, err := GetRandomBytes(length)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// GetRandomBytes retrieves a byte array of 'length'
func GetRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

// GetHashingAlgorithName retrieves a string representation of the hashing algorithm
func GetHashingAlgorithmName(alg crypto.Hash) string {
	switch alg {
	case crypto.SHA1:
		return "SHA1"
	case crypto.SHA256:
		return "SHA-256"
	case crypto.SHA384:
		return "SHA-384"
	case crypto.SHA512:
		return "SHA-512"
	}
	return ""
}

// GetHash returns a byte array to the hash of the data.
// alg indicates the hashing algorithm. Currently, the only supported hashing algorithms
// are SHA1, SHA256, SHA384 and SHA512
func GetHashData(data []byte, alg crypto.Hash) ([]byte, error) {

	if data == nil {
		return nil, fmt.Errorf("Error - data pointer is nil")
	}

	switch alg {
	case crypto.SHA1:
		s := sha1.Sum(data)
		return s[:], nil
	case crypto.SHA256:
		s := sha256.Sum256(data)
		return s[:], nil
	case crypto.SHA384:
		//SHA384 is implemented in the sha512 package
		s := sha512.Sum384(data)
		return s[:], nil
	case crypto.SHA512:
		s := sha512.Sum512(data)
		return s[:], nil
	}

	return nil, fmt.Errorf("Error - Unsupported hashing function %d requested. Only SHA1, SHA256, SHA384 and SHA512 supported", alg)
}

const certSubjectName = "ISecl Test Cert"
const certExpiryDays = 180

// CreateTestCertAndRSAPrivKey is a helper function to create a test certificate. This will be primarily used
// by test function to make a keypair and certificate that can be used for encryption and signing.
func CreateTestCertAndRSAPrivKey(bits ...int) (*rsa.PrivateKey, string, error) {
	if len(bits) > 1 {
		return nil, "", fmt.Errorf("Error: Function only accepts a single parameter - RSA keylength in bits - 1024 or 2048")
	}
	var rsaBitLength int

	if len(bits) == 0 {
		rsaBitLength = 2048
	} else {
		rsaBitLength = bits[0]
	}

	if !(rsaBitLength == 2048 || rsaBitLength == 1024) {
		return nil, "", fmt.Errorf("Error: RSA keylength support is only 1024 or 2048")
	}

	rsaKeyPair, err := rsa.GenerateKey(rand.Reader, rsaBitLength)
	if err != nil {
		return nil, "", err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{certSubjectName},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * certExpiryDays),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &rsaKeyPair.PublicKey, rsaKeyPair)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}

	out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	return rsaKeyPair, out.String(), nil

}

func HashAndSignPKCS1v15(data []byte, rsaPriv *rsa.PrivateKey, alg crypto.Hash) ([]byte, error) {

	hash, err := GetHashData(data, alg)
	if err != nil {
		return nil, err
	}
	return rsa.SignPKCS1v15(rand.Reader, rsaPriv, alg, hash)

}
