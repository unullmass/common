package model

// EncryptionHeader is used append to the encrypted file during AES GCM encryption mode
type EncryptionHeader struct {
	MagicText            [12]byte
	OffsetInLittleEndian uint32
	Version              [4]byte
	IV                   [12]byte
	EncryptionAlgorithm  [12]byte
}
