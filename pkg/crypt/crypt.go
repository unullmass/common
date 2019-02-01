package crypt

// SignedData contains the original byte stream that was signed along with
type SignedData struct {
	Data      []byte `json:"data"` //json formatted VMtrust report.
	Alg       string `json:"hash_alg"`
	Cert      string `json:"cert"` //pem formatted certificate
	Signature []byte `json:"signature"`
}
