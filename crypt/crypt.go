package crypt

import (
	"crypto"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
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
