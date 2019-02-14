package os

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestIsFileEncryptedTrue(t *testing.T) {
	encImagePath := "../test/cirros-x86.qcow2_enc"
	isImageEncrypted, err := IsFileEncrypted(encImagePath)
	assert.NoError(t, err)
	assert.True(t, isImageEncrypted)
}

func TestIsFileEncryptedFalse(t *testing.T) {
	encImagePath := "../test/cirros-x86.qcow2"
	isImageEncrypted, err := IsFileEncrypted(encImagePath)
	assert.NoError(t, err)
	assert.False(t, isImageEncrypted)
}
