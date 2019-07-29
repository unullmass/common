/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
 package setup

 import (
		 "fmt"
		 "flag"
		 "io"
		 "os"
		 "crypto/x509"
		 "bytes"
		 "strings"
		 "errors"
		 "io/ioutil"
		 "net/http"
		 "net/url"
		 "crypto/tls"
		 "encoding/pem"
		 "intel/isecl/lib/common/validation"
		 "intel/isecl/lib/common/crypt"
		 cos "intel/isecl/lib/common/os"
 )
 
 type Download_Cert struct {
		 Flags              []string
		 KeyFile            string     
		 CertFile           string 
		 KeyAlgorithm       string
		 KeyAlgorithmLength int
		 CommonName         string
		 SanList            string
		 CertType           string
		 CaCertsDir         string
		 BearerToken        string
	     ConsoleWriter      io.Writer
 }

 func GetCertificateFromCMS(certType string, keyAlg string, keyLen int, cmsBaseUrl string, commonName string, hosts string, caCertsDir string, bearerToken string) (key []byte, cert []byte, err error) {
	 //TODO: use CertType for TLS or Signing cert	
	csrData, key, err := crypt.CreateKeyPairAndCertificateRequest(commonName, hosts, keyAlg, keyLen)
	if err != nil {
	   return nil, nil, fmt.Errorf("Certificate setup: %v", err)
	}	

   url, err := url.Parse(cmsBaseUrl)
   if err != nil {
		   fmt.Println("Configured CMS URL is malformed: ", err)
		   return nil, nil, fmt.Errorf("Certificate setup: %v", err)
   }
   certificates, _ := url.Parse("certificates?certType=" + certType)
   endpoint := url.ResolveReference(certificates)
   csrPemBytes := pem.EncodeToMemory(&pem.Block{Type: "BEGIN CERTIFICATE REQUEST", Bytes: csrData})
   req, err := http.NewRequest("POST", endpoint.String(),  bytes.NewBuffer(csrPemBytes))
   if err != nil {
		   fmt.Println("Failed to instantiate http request to CMS")
		   return nil, nil, fmt.Errorf("Certificate setup: %v", err)
   }
   req.Header.Set("Accept", "application/x-pem-file")        
   req.Header.Set("Content-Type", "application/x-pem-file")      
   req.Header.Set("Authorization", "Bearer " + bearerToken)  
   
   rootCaCertPems, err := cos.GetDirFileContents(caCertsDir, "*.pem" )
	if err != nil {
		return nil, nil, err
	}

    rootCAs, _ := x509.SystemCertPool()
    if rootCAs == nil {
		rootCAs = x509.NewCertPool()
    }
	for _, rootCACert := range rootCaCertPems{
		if ok := rootCAs.AppendCertsFromPEM(rootCACert); !ok {
	                return nil, nil, err
	        }
	}

   client := &http.Client{
		   Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: false,
					RootCAs: rootCAs,
				},
		   },
   }
   resp, err := client.Do(req)
   if err != nil {
		   fmt.Println("Failed to perform HTTP request to CMS")
		   return nil, nil, fmt.Errorf("Certificate setup: %v", err)
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
		   text, _ := ioutil.ReadAll(resp.Body)
		   errStr := fmt.Sprintf("CMS request failed to download Certificate (HTTP Status Code: %d)\nMessage: %s", resp.StatusCode, string(text))
		   fmt.Println(errStr)
		   return nil, nil, fmt.Errorf("Certificate setup: %v", err)
   }
   cert, err = ioutil.ReadAll(resp.Body)
   if err != nil {
		   fmt.Println("Failed to read CMS response body")
		   return nil, nil, fmt.Errorf("Certificate setup: %v", err)
   }   
	return
}
 
 func (tc Download_Cert) Run(c Context) error {
		 fmt.Fprintln(tc.ConsoleWriter, "Running Certificate download setup...")
		 fs := flag.NewFlagSet("download_cert", flag.ContinueOnError)
		 force := fs.Bool("force", false, "force recreation, will overwrite any existing certificate")
		 certType := fs.String("cert", tc.CertType, "type of the certificate")
		 
		 err := fs.Parse(tc.Flags)
		 if err != nil {				 
				 return errors.New("Certificate setup: Unable to parse flags") 
		 }
		 fmt.Println("Certificate Type :" + *certType)
		 cmsBaseUrl, err := c.GetenvString("CMS_BASE_URL", "CMS base URL in https://{{cms}}:{{cms_port}}/cms/v1/")
	     if err != nil || cmsBaseUrl == "" {			     
				 return errors.New("Certificate setup: CMS_BASE_URL not found in environment for Download Certificate") 
		 }

		keyFile, err := c.GetenvString("KEY_PATH", "Path of file where key needs to be stored")
		if err == nil {
			tc.KeyFile = keyFile
		}

		certPath, err := c.GetenvString("CERT_PATH", "Path of file/directory where certificate needs to be stored")
		if err == nil {
			tc.CertFile = certPath
		}

		commonName, err := c.GetenvString("COMMON_NAME", "Common Name of certificate")
		if err == nil {
			tc.CommonName = commonName
		}

		defaultHostname, err := c.GetenvString("SAN_LIST", "Comma separated list of hostnames to add to Certificate")
		if err != nil {
			defaultHostname = tc.SanList
		}
		host := fs.String("host_names", defaultHostname, "Comma separated list of hostnames to add to Certificate")
 
		bearerToken := tc.BearerToken
		tokenFromEnv, err := c.GetenvString("BEARER_TOKEN", "bearer token")
	    if err == nil {
			bearerToken = tokenFromEnv
		}
		if bearerToken == "" {			
			return errors.New("Certificate setup: BEARER_TOKEN not found in environment for Download Certificate") 
		}
 
		 if *force || tc.Validate(c) != nil {
			if *host == "" {
				return errors.New("Certificate setup: no SAN hostnames specified")
			}
			hosts := strings.Split(*host, ",")
	
			// validate host names
			for _, h := range hosts {
				valid_err := validation.ValidateHostname(h)
				if valid_err != nil {
					return valid_err
				}
			}
			key, cert, err := GetCertificateFromCMS(tc.CertType, tc.KeyAlgorithm, tc.KeyAlgorithmLength, cmsBaseUrl, tc.CommonName, *host, tc.CaCertsDir, bearerToken)
			if err != nil {
				return fmt.Errorf("Certificate setup: %v", err)
			}
			err = crypt.SavePrivateKeyAsPKCS8(key, tc.KeyFile)
			if err != nil {
				return fmt.Errorf("Certificate setup: %v", err)
			} 

			fi, err := os.Stat(tc.CertFile)
			if err != nil || fi.Mode().IsRegular() {
				err = ioutil.WriteFile(tc.CertFile, cert, 0644)
				if err != nil {
					fmt.Println("Could not store Certificate")
					return fmt.Errorf("Certificate setup: %v", err)
				}				
				os.Chmod(tc.CertFile, 0644)
			} else if fi.Mode().IsDir() {
				err = crypt.SavePemCertWithShortSha1FileName(cert, tc.CertFile)
				if err != nil {
					fmt.Println("Could not store Certificate")
					return fmt.Errorf("Certificate setup: %v", err)
				}
			}
		 } else {
				 fmt.Println("Certificate already downloaded, skipping")
		 }           
		  return nil  
 }
 
 func (tc Download_Cert) Validate(c Context) error {	 
	fmt.Fprintln(tc.ConsoleWriter, "Validating Certificate download setup...")	

	 _, err := os.Stat(tc.KeyFile)
	 if os.IsNotExist(err) {
		 return errors.New("KeyFile is not configured")
	 }
	 return nil
  }