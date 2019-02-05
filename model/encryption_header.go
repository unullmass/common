package model

// EncryptionHeader is used append to the encrypted file during AES GCM encryption mode
type EncryptionHeader struct {
	MagicText [12]byte
	Offset    uint32
	IV        [12]byte
	Algorithm [6]byte
}
