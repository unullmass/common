package crypt

import (
	"io/ioutil"
	"strings"
)

// EncryptionHeader is used append to the encrypted file during AES GCM encryption mode
type EncryptionHeader struct {
	MagicText            [12]byte
	OffsetInLittleEndian uint32
	Version              [4]byte
	IV                   [12]byte
	EncryptionAlgorithm  [12]byte
}

// EncryptionHeaderExists method is used to check if the file is encryped and returns a boolean value.
// TODO : move it a different package where all the ISecL specific functions are added
func EncryptionHeaderExists(encFilePath string) (bool, error) {

	var encryptionHeader EncryptionHeader
	//check if file is encrypted
	encFileContents, err := ioutil.ReadFile(encFilePath)
	if err != nil {
		return false, err
	}

	magicText := encFileContents[:len(encryptionHeader.MagicText)]
	if !strings.Contains(string(magicText), EncryptionHeaderMagicText) {
		return false, nil
	}
	return true, nil
}
