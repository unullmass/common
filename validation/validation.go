/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package validation

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
)

// var uname_reg, _ = regexp.Compile(`^[A-Za-z]{1}[A-Za-z0-9_]{1,31}$`)

// var pwd_reg, _ = regexp.Compile("^[A-Za-z0-9#?!@$%^&*-._]{1,}$")
// var hostname_reg, _ = regexp.Compile("^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$")
// var ip_reg, _ = regexp.Compile("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$")
// var idf_reg, _ = regexp.Compile(`^[a-zA-Z_]{1}[a-zA-Z0-9_]{1,127}$`)

// var hostname_reg, _ = regexp.Compile(`^(?!\s*$).+`)
// var ip_reg, _ = regexp.Compile(`^(?!\s*$).+`)

var (
	unameReg    = regexp.MustCompile(`^[A-Za-z]{1}[A-Za-z0-9_]{1,31}$`)
	hostnameReg = regexp.MustCompile(`.+`)
	ipReg       = regexp.MustCompile(`.+`)
	idfReg      = regexp.MustCompile(`^[a-zA-Z_]{1}[a-zA-Z0-9_]{1,127}$`)
	hardwareuuid_reg = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
)

// ValidateEnvList can check if all environment variables in input slice exist
// If things missing, return a slice contains all missing variables and an error.
// Otherwise two nil
func ValidateEnvList(required []string) ([]string, error) {

	missing := make([]string, 0)
	for _, e := range required {
		if _, exist := os.LookupEnv(e); !exist {
			missing = append(missing, e)
		}
	}
	if len(missing) == 0 {
		return nil, nil
	} else {
		return missing, errors.New("One or more required environment variables are missing")
	}
}

// ValidateURL checks if input URL meets the constraints
// protocols is a set of protocols the URL can be supporting, mimicked by a map
// protocol names will be stored in key, and the value is discard
// path is the specified path the URL must follow
// Returns an error if any requirement is not met
func ValidateURL(testURL string, protocols map[string]byte, path string) error {

	urlObj, err := url.Parse(testURL)
	if err != nil {
		return errors.New("Invalid base URL")
	}
	if _, exist := protocols[urlObj.Scheme]; !exist {
		return errors.New("Unsupported protocol")
	}
	if urlObj.Path != path {
		return errors.New("Invalid path in URL")
	}
	if urlObj.RawQuery != "" || urlObj.Fragment != "" {
		return errors.New("Unexpected inputs")
	}
	return nil
}

// ValidateAccount information, both username and password primarily forbidding any quotation marks in the input
// as well as restricting the length for both of them
func ValidateAccount(uname string, pwd string) error {

	if unameReg.MatchString(uname) && pwd != "" {
		return nil
	}
	return errors.New("Invalid input for username or password")
}

// ValidateHostname method is used to validate the hostname string
func ValidateHostname(hostname string) error {

	if hostnameReg.MatchString(hostname) || ipReg.MatchString(hostname) {
		return nil
	}
	return errors.New("Invalid hostname or ip")
}

// ValidateInteger method is used to validate an input integer value
func ValidateInteger(number string, cnt int) error {

	mat, err := regexp.Match(fmt.Sprintf("^[0-9]{1,%d}$", cnt), []byte(number))
	if err != nil {
		return err
	}
	if !mat {
		return errors.New("Invalid numeric string")
	}
	return nil
}

// ValidateRestrictedString method is used to validate a string based on allowed characters
func ValidateRestrictedString(str string, allowed string) error {

	mat, err := regexp.Match(fmt.Sprintf("^[%s]{1,128}$", allowed), []byte(str))
	if err != nil {
		return err
	}
	if !mat {
		return errors.New("Invalid input string")
	}
	return nil
}

// ValidateIdentifier method is used to validate an identifier value
func ValidateIdentifier(idf string) error {

	if idfReg.MatchString(idf) {
		return nil
	}
	return errors.New("Invalid identifier")
}

// ValidateStrings method is used to validate input strings
func ValidateStrings(strings []string) error {
	fmt.Println("inside the method")
	strRegEx, err := regexp.Compile("(^[a-zA-Z0-9_///.-]*$)")
	if err != nil {
		fmt.Println("Error occured")
		return err
	}

	for _, stringValue := range strings {
		fmt.Println("String value: ", stringValue)
		if !strRegEx.MatchString(stringValue) {
			return errors.New("Invalid string formatted input")
		}
	}
	return nil
}

// ValidatePemEncodedKey method is used to validate input keysin string format
func ValidatePemEncodedKey(key string) error {
	strRegEx, err := regexp.Compile("(^[-a-zA-Z0-9//=+ ]*$)")
	if err != nil {
		return err
	}

	if !strRegEx.MatchString(key) {
		return errors.New("Invalid key format")
	}
	return nil
}

// ValidateBase64String method checks if a string has a valid base64 format
func ValidateBase64String(value string) error {
	r := regexp.MustCompile("^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$")
	if !r.MatchString(value) {
		return errors.New("Invalid digest format")
	}
	return nil
}

// ValidateHexString method checks if a string has a valid hex format
func ValidateHexString(value string) error {
	r := regexp.MustCompile("^[a-fA-F0-9]+$")
	if !r.MatchString(value) {
		return errors.New("Invalid hex string format")
	}
	return nil
}

// ValidateUUID method is used to check if the given UUID is of valid format
func ValidateUUID(uuid string) error {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	if !r.MatchString(uuid) {
		return errors.New("Invalid UUID format")
	}
	return nil
}

func ValidateHardwareUUID(uuid string) error {
	if hardwareuuid_reg.MatchString(uuid) {
		return nil
	}
	return errors.New("Invalid hardware uuid")
}
