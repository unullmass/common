/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package validation

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEnvList(t *testing.T) {
	a := assert.New(t)

	os.Setenv("GOENV_TEST1", "NOT IMPORTATN CONTENT")
	os.Setenv("GOENV_TEST2", "NOT IMPORTATN CONTENT")
	os.Setenv("GOENV_TEST3", "NOT IMPORTATN CONTENT")

	neededEnv := []string{"GOENV_TEST1", "GOENV_TEST2", "GOENV_TEST3"}
	missing, err := ValidateEnvList(neededEnv)

	a.Equal(nil, err)
	a.Equal(0, len(missing))

	moreThanNeededEnv := []string{"GOENV_TEST2", "GOENV_TEST3", "GOENV_TEST4", "GOENV_TEST5"}

	missing1, err1 := ValidateEnvList(moreThanNeededEnv)
	sampleErr := errors.New("One or more required environment variables are missing")

	a.Equal(sampleErr, err1)

	// should show [GOENV_TEST4:0 GOENV_TEST5:0]
	// run with with go test -v
	fmt.Println(missing1)
}

func TestValidateURL(t *testing.T) {
	a := assert.New(t)

	goodURL1 := "https://google.com"
	goodURL2 := "http://good.url.with.port:5566"
	goodURL3 := "https://good.url.https.with.port:5566"

	goodTests := []string{goodURL1, goodURL2, goodURL3}

	protocols := make(map[string]byte)
	protocols["http"] = 0
	protocols["https"] = 0

	for _, goodStr := range goodTests {
		err := ValidateURL(goodStr+"/tds/", protocols, "/tds/")
		a.Equal(nil, err)
	}

	goodURL4 := "https://10.1.68.68:443/v1/keys/438c7486-9827-4072-89c5-93ae4538114e/transfer"
	err := ValidateURL(goodURL4, protocols, "/v1/keys/438c7486-9827-4072-89c5-93ae4538114e/transfer")
	a.Equal(nil, err)

	badURL1 := "bad.url.without.protocol/tds/"
	badURL2 := "scheme://bad.url.with.wrong.protocol/tds/"
	badURL3 := "https://bad.url.with.path/tds/path/path/"
	badURL4 := "https://bad.url.with.query/tds/?query=haha"

	// err_1 := errors.New("Invalid base URL")
	err2 := errors.New("Unsupported protocol")
	err3 := errors.New("Invalid path in URL")
	err4 := errors.New("Unexpected inputs")

	badTests := []string{badURL1, badURL2, badURL3, badURL4}
	badResult := []error{err2, err2, err3, err4}

	for i, badStr := range badTests {

		err := ValidateURL(badStr, protocols, "/tds/")
		a.Equal(badResult[i], err)
	}
}

func TestValidateAccount(t *testing.T) {
	// a := assert.New(t)

	// good_uname_1 := "uname.has-symbol"
	// good_uname_2 := "abcd1234"
	// good_uname_3 := "a1a2_3d4f5g6_7j8k9l1_2s3d45g6h7"

	// good_pwd_1 := "easy_guess123"
	// good_pwd_2 := "hardone_A-Za-z0-9#?!@$%^&*-"
	// good_pwd_3 := "tooshort"

	// good_uname_tests := []string{good_uname_1, good_uname_2, good_uname_3}
	// good_pwd_tests := []string{good_pwd_1, good_pwd_2, good_pwd_3}

	// for _, uname := range good_uname_tests {

	// 	for _, pwd := range good_pwd_tests {
	// 		// fmt.Println(uname + " " + pwd)
	// 		err := ValidateAccount(uname, pwd)
	// 		a.Equal(nil, err)
	// 	}
	// }

	// bad_uname_1 := "fishy_symbols \" ` '"
	// bad_uname_2 := "\" \" \" UNION SELECT * FROM"
	// bad_uname_3 := "))) OR TRUE"
	// bad_uname_4 := "12_number_start_with"

	// bad_pwd_1 := "fishy_symbols \" ` '"
	// bad_pwd_2 := "\" \" \" UNION SELECT * FROM"
	// bad_pwd_3 := "))) OR TRUE"
	// bad_pwd_4 := "tooshrt"
	// bad_pwd_5 := ""

	// bad_ret := errors.New("Invalid input for username or password")

	// bad_uname_tests := []string{bad_uname_1, bad_uname_2, bad_uname_3, bad_uname_4}
	// bad_pwd_tests := []string{bad_pwd_1, bad_pwd_2, bad_pwd_3, bad_pwd_4, bad_pwd_5}

	// for _, bad_uname := range bad_uname_tests {
	// 	fmt.Println(bad_uname + " " + good_pwd_1)
	// 	err := ValidateAccount(bad_uname, good_pwd_1)
	// 	a.Equal(bad_ret, err)
	// }
	// for _, bad_pwd := range bad_pwd_tests {
	// 	err := ValidateAccount(good_uname_1, bad_pwd)
	// 	a.Equal(bad_ret, err)
	// }
}

func TestValidateStrings(t *testing.T) {
	goodString1 := "/var/lib/nova/instances/"
	goodString2 := "abcd1234"
	goodString3 := "workload-agent"
	goodString4 := "cirros-x86.qcow2_enc"
	goodString5 := ""
	goodString6 := "cirros-enc.qcow2"

	goodStringValueArr := []string{goodString1, goodString2, goodString3, goodString4, goodString5, goodString6}
	err := ValidateStrings(goodStringValueArr)
	assert.NoError(t, err)

	badString1 := "fishy_symbols \" ` ' )))((*&^ "
	badString2 := "\" \" \" SELECT * FROM TABLE"
	badString3 := "<inputTag>"

	badStringValueArr := []string{badString1, badString2, badString3}
	err = ValidateStrings(badStringValueArr)
	assert.Error(t, err)
}

func TestValidatePemEncodedKey(t *testing.T) {

	goodString1 := "-----BEGIN CERTIFICATE----MIIEoDCCA4igAwIBAgIIHZR9rDPTS9IwDQYJKoZIhvcNAQELBQAwGzEZMBcGA1UEAxMQbXR3aWxzb24tcGNhLWFpazAeFw0xOTAzMDgwOTQ0MDJaFw0yOTAzMDUwOTQ0MDJaMCUxIzAhBgNVBAMMGkNOPUJpbmRpbmdfS2V5X0NlcnRpZmljYXRlMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvarJulWowk4Mr0CwJ8SMrGNvbHKryZmTgLFftnj7bwd7VUx0zACjwv0jObFAonocBMEO39hxhHeQWuyUB1OEuvt3Kzg40EG11PvsHttWM2rgYk88z+E7vsQrdHx08FLeN4T+SK9ML+uKDPFuQWLrzZ2irQxGznokn2aAz+zl8vIDkVjmVDw4I9D6/6hMAEaE4WGiznoPS1sEe6vyqPhOL9dQOrFzW5qC6JKNwPtuYIWg5VrhrxguTHe5vK7TpZxgbxjWbeSdiDAnNU7UdZM5QboJ8Ao7Oz54OtIqpUuIEXSzIaJ/Dbo0yXrYCJzcSDTpSDimitqfTiot9p6+pQSM8wIDAQABo4IB3DCCAdgwDgYDVR0PAQH/BAQDAgUgMIGdBgdVBIEFAwIpBIGR/1RDR4AXACIAC0y1Y5XeuPW1DaT1Y11paGPfMQ3JGXM61ail1pp716pUAAQA/1WqAAAAAAgwHmMAAAAGAAAAAQEABwA+AAw2AAAiAAvKQw0fWfgw6nQE4MTsGY7tTwepogu1YED+3ZSs+rkBbgAiAAvZHMHbKCiUnL7jxkTMqhbMNtoyKqi7eg5Tm7eQFQM3+jCCARQGCFUEgQUDAikBBIIBBgAUAAsBADiGza+SuXvKAkfBWbYLXXmneXfAUaUtynArnhey7icOvBNqNVVROXx9iaH5wxgmzshJRDXIkJhYpZ8xfNkFqfPkAtHSw3IugrTe4bokzVwVMe+a0c+OT9ptDEkb7TiUKNc0hX2O0jc4dPSzYBdAn/3HSnV1f5DQBKT/QTk/Io+7F0Vu2phz96Ax7d73mjufehE920hq77v/mVXNOtmyQ6Q9OXQCgDwgfAQwBc6iuBLsXPAK1GUN2D/8N29+CQQfa0/KzHrHc/Ioh+PTOgg7TaYmJhB+NasdzV7t3YlAa0x6ZPPZwagpSnbrGpnIw2H0Wvy/YebfvGskymvvsiSwUGQwDgYIVQSBBQMCKQIEAgAAMA0GCSqGSIb3DQEBCwUAA4IBAQBh6LiD+zec5Tp0qZnpNniw519/JHOVN3HDcF1mv/BKSEeu7BmWp62c3Agf6+8anWsrpTg936DAgCKUgjJN2+m5fdJhNOUs662/PlLE17FjeAYMIxfNlHlP5nNRq5F6T8I1BDaCsT4dgmQlqUkCrHRaPB9nqajPOjYihYpSxNnUT/b4NFnK8T7gEvrfGM1EF6m1dPx4IQSxznBD4XwjBU/KjjVEuPFks/bPwf/sz5t3itnpJTFazrnBp4Wr2kzxKLoPDKmhtgStE0mArpLWnEbgMXGmkm5mgzHCFs5phEghhkAHfyT/lLjMKzzkN7BmjL6qM9TIA67g8LIMuI/wyKXe-----END CERTIFICATE-----"
	goodString2 := "tNPEL9e3O97nQ2nTlzwyPxxNLAjwwSvTnGcVeOmIfbboqBZkixDx3TR6/n13s7eury6LVqOjEOsItExy/BvXE+g6WTP61dEnenM1fNlmNY/O6bNqEwgxVEfamdNBN9FYuB2inzwueEiB5J2WtkCB+0c2pr86oMl8eYlUve1r0+VRbOhQVhMKskVLKZonBxvLrup1geqIDD8DuK/EVt8aCYfjY5n0B1NtfSwQpNzWjmTc0g9H0aH/tYerfB6Q+KqVo2be1wWQscNBfVvVGjwB3brT6k2ODUSCDlGJV5IhVNq+ooM5s+o21c8+fX1bXiFOdA7T/TnuCKJyEMxrAvzypQ=="
	
	err := ValidatePemEncodedKey(goodString1)
	assert.NoError(t, err)

	err = ValidatePemEncodedKey(goodString2)
	assert.NoError(t, err)

	badString1 := "fishy_symbols \" ` ' )))((*&^ "
	badString2 := "\" \" \" SELECT * FROM TABLE"
	badString3 := "<inputTag>"

	err = ValidatePemEncodedKey(badString1)
	assert.Error(t, err)

	err = ValidatePemEncodedKey(badString2)
	assert.Error(t, err)
	
	err = ValidatePemEncodedKey(badString3)
	assert.Error(t, err)
}

func TestValidateUUIDv4(t *testing.T) {

	goodUUID := "2a16f0bd-aa32-458d-9d7b-fe1a7048c3e2"	
	err := ValidateUUIDv4(goodUUID)
	assert.NoError(t, err)

	badUUID1 := "sample-string-with-novalidValue"
	badUUID2 := "$%^&#$--1234-7890-2345"
	badUUID3 := "fChs4vXGYJ6hUHCILkLNg1STbF3YC270tVb3pUP9AfA="

	err = ValidateUUIDv4(badUUID1)
	assert.Error(t, err)

	err = ValidateUUIDv4(badUUID2)
	assert.Error(t, err)
	
	err = ValidateUUIDv4(badUUID3)
	assert.Error(t, err)
}

func TestValidateBase64String(t *testing.T) {

	goodUUID := "fChs4vXGYJ6hUHCILkLNg1STbF3YC270tVb3pUP9AfA="	
	err := ValidateBase64String(goodUUID)
	assert.NoError(t, err)

	badUUID1 := "sample-string-with-novalidValue"
	badUUID2 := "$%^&#$--1234-7890-2345"
	badUUID3 := "2a16f0bd-aa32-458d-9d7b-fe1a7048c3e2"

	err = ValidateBase64String(badUUID1)
	assert.Error(t, err)

	err = ValidateBase64String(badUUID2)
	assert.Error(t, err)
	
	err = ValidateBase64String(badUUID3)
	assert.Error(t, err)
}

func TestValidateHexString(t *testing.T) {

	goodUUID := "f41d0ebcaafb14a3a9a16329bd4c1248db2759576e1e8992c342f6cb52bab333"	
	err := ValidateHexString(goodUUID)
	assert.NoError(t, err)

	badUUID1 := "sample-string-with-novalidValue"
	badUUID2 := "$%^&#$--1234-7890-2345"
	badUUID3 := "2a16f0bd-aa32-458d-9d7b-fe1a7048c3e2"

	err = ValidateHexString(badUUID1)
	assert.Error(t, err)

	err = ValidateHexString(badUUID2)
	assert.Error(t, err)
	
	err = ValidateHexString(badUUID3)
	assert.Error(t, err)
}

func TestValidateXMLString(t *testing.T) {

	goodXML := "<domain type='kvm'><name>instance-00000035</name><uuid>6f62d266-7dfd-4eaa-b31c-fc27a68565a9</uuid>"+
               "<metadata><nova:instance xmlns:nova=\"http://openstack.org/xmlns/libvirt/nova/1.0\"><nova:package version=\"16.0.1\"/>" +
                "<nova:name>jy</nova:name><nova:creationTime>2018-06-14 18:28:00</nova:creationTime><nova:flavor name=\"small\">"+
				"<nova:memory>4000</nova:memory><nova:disk>1</nova:disk><nova:swap>0</nova:swap><nova:ephemeral>1</nova:ephemeral>" +
				"<nova:vcpus>1</nova:vcpus></nova:flavor></metadata></domain>"

	err := ValidateXMLString(goodXML)
	assert.NoError(t, err)

	badXML := "$%^&#$--1234-7890-2345"
	err = ValidateXMLString(badXML)
	assert.Error(t, err)
}

func TestValidateDate(t *testing.T) {

	goodDate1 := "31-07-2010"

	err := ValidateDate(goodDate1)
	assert.NoError(t, err)

	badDate1 := "1/13/2010"
	badDate2 := "29-02-200a"

	err = ValidateDate(badDate1)
	assert.Error(t, err)

	err = ValidateDate(badDate2)
	assert.Error(t, err)
}
