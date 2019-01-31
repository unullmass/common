package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os/exec"
)

// ExecuteCommand is used to execute a linux command line command and return the output of the command with an error if it exists.
func ExecuteCommand(cmd string, args []string) (string, error) {
	out, err := exec.Command(cmd, args...).Output()
	return string(out), err
}

// ExecuteCommandCentext is used to execute a linux command line command with context provided
// and return the output of the command with an error if it exists.
func ExecuteCommandCentext(ctx context.Context, cmd string, args []string) (string, error) {
	out, err := exec.CommandContext(ctx, cmd, args...).Output()
	return string(out), err
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
