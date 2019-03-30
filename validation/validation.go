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
	uname_reg    = regexp.MustCompile(`^[A-Za-z]{1}[A-Za-z0-9_]{1,31}$`)
	hostname_reg = regexp.MustCompile(`.+`)
	ip_reg       = regexp.MustCompile(`.+`)
	idf_reg      = regexp.MustCompile(`^[a-zA-Z_]{1}[a-zA-Z0-9_]{1,127}$`)
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
func ValidateURL(test_url string, protocols map[string]byte, path string) error {

	url_obj, err := url.Parse(test_url)
	if err != nil {
		return errors.New("Invalid base URL")
	}
	if _, exist := protocols[url_obj.Scheme]; !exist {
		return errors.New("Unsupported protocol")
	}
	if url_obj.Path != path {
		return errors.New("Invalid path in URL")
	}
	if url_obj.RawQuery != "" || url_obj.Fragment != "" {
		return errors.New("Unexpected inputs")
	}
	return nil
}

// Validate account information, both user name and pass word
// Primarily forbidding any quotation marks in the input
// As well as restricting the length for both of them
func ValidateAccount(uname string, pwd string) error {

	if uname_reg.MatchString(uname) && pwd != "" {
		return nil
	}
	return errors.New("Invalid input for username or password")
}

func ValidateHostname(hostname string) error {

	if hostname_reg.MatchString(hostname) || ip_reg.MatchString(hostname) {
		return nil
	}
	return errors.New("Invalid hostname or ip")
}

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

func ValidateIdentifier(idf string) error {

	if idf_reg.MatchString(idf) {
		return nil
	}
	return errors.New("Invalid identifier")
}
