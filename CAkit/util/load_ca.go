// load_ca.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>
package util


import (
 	"crypto/rsa"
	"crypto/x509"
 	"encoding/pem"
 	"errors"
 	"io/ioutil"

	cakit "armador/CAkit"
)


func LoadCA(certFName, keyFName string, callback cakit.PasswordCallback, kdf cakit.KDF) (cert *x509.Certificate, privkey *rsa.PrivateKey, the_err error) {
	
	cert	=nil
	privkey =nil
	the_err =nil
		
	// ** Load CA
	// Can't use LoadX509KeyPair since KEY should be encrypted ...
	
	// * Load CA cert first
	certPEM, err := ioutil.ReadFile(certFName)
	if nil != err {
		the_err=errors.New("Could not read PEM")
 		return
	}
	// ioutil auto-closes file
	
	certPEMBlock,_ := pem.Decode(certPEM)
	if "CERTIFICATE" != certPEMBlock.Type {
		the_err=errors.New("Load CA Cert: "+err.Error())
		return
	}
	
	// * Load (potentially encrypted) KEY from PEM
	keyPEM, err := ioutil.ReadFile(keyFName+".key")
	if nil != err  {
		the_err=errors.New("Load CA Key: "+err.Error())
		return
	}
	// ioutil auto-closes file
	
	// Decode PEM-encoded key into a block ...
	keyPEMBlock,_ := pem.Decode(keyPEM)
	
	if "RSA PRIVATE KEY" != keyPEMBlock.Type {
		the_err = errors.New("Malformed cert")
		return
	}
	
	var keyBytes []byte
	
	// IsEncryptedPEMBlock(b *pem.Block) bool
	if x509.IsEncryptedPEMBlock(keyPEMBlock) {
		
		var key []byte
		
		// Assume we encrypted it, so we know the params ;)
		password, err := callback()
		if nil!=err {
			the_err = err
			return
		}
		
		if nil!=kdf {
			if key,err = kdf(password); nil!=err {
				the_err = err
				return
			}
		} else {
			key=[]byte(password)
		}
		
		keyBytes,err = x509.DecryptPEMBlock(keyPEMBlock, key)
		if nil!=err {
			the_err = err
			return
		}
	} else {
		keyBytes = keyPEMBlock.Bytes
	}

	certdata,err := x509.ParseCertificate(certPEMBlock.Bytes)
	if nil!=err {
		the_err = err
		return
	}
	
	pk,err := x509.ParsePKCS1PrivateKey(keyBytes)
	if nil!=err {
		the_err = err
		return
	}
	
	return certdata,pk,nil
}
