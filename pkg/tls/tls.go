package tls

import (
	"crypto/sha256"
	"crypto/x509"
	"errors"
)

func verifyCertBySha256(certSha256 [32]byte) func([][]byte, [][]*x509.Certificate) error {
	return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		if len(rawCerts) <= 0 {
			return errors.New("kms-client tls: no certificates supplied")
		}
		hostRawCert := rawCerts[0]
		fingerprint := sha256.Sum256(hostRawCert)
		if fingerprint != certSha256 {
			return errors.New("kms-client tls: fingerprint does not match")
		}
		hostCert, err := x509.ParseCertificate(hostRawCert)
		if err != nil {
			return errors.New("kms-client tls: could not parse certificate")
		}
		intermediates := x509.NewCertPool()
		roots, err := x509.SystemCertPool()
		if err != nil {
			roots = x509.NewCertPool()
		}
		if hostCert.IsCA {
			// we have a single, self signed certificate
			roots.AddCert(hostCert)
		}
		rest := rawCerts[1:]
		for _, rawCert := range rest {
			cert, err := x509.ParseCertificate(rawCert)
			if err != nil {
				return errors.New("kms-client tls: failed to parse x509 certificate")
			}
			if cert.IsCA {
				roots.AddCert(cert)
			} else {
				intermediates.AddCert(cert)
			}
		}
		opts := x509.VerifyOptions{
			Intermediates: intermediates,
			Roots:         roots,
		}
		_, err = hostCert.Verify(opts)
		return err
	}
}
