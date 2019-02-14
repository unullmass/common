/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package os

import (
	"os"
	"path/filepath"
)

// ChownR method is used to change the ownership of all the file in a directory
func ChownR(path string, uid, gid int) error {
	return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if err == nil {
			err = os.Chown(name, uid, gid)
		}
		return err
	})
}

// IsFileEncrypted method is used to check if the file is encryped and returns a boolean value.
// TODO : move it a different package where all the ISecL specific functions are added
func IsFileEncrypted(encFilePath string) (bool, error) {

	var encryptionHeader crypt.EncryptionHeader
	//check if file is encrypted
	encFileContents, err := ioutil.ReadFile(encFilePath)
	if err != nil {
		return false, err
	}

	magicText := encFileContents[:len(encryptionHeader.MagicText)]
	if !strings.Contains(string(magicText), crypt.EncryptionHeaderMagicText) {
		return false, nil
	}
	return true, nil
}
