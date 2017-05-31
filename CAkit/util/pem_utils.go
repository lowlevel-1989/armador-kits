// pem_utils.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>
package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
// 	"crypto/x509/pkix"
	"encoding/pem"
 	"errors"
// 	"fmt"
	"io/ioutil"

	"os"
	
	cakit ".."
)

const (
	
	keyType		= "RSA"
	blockCert	= "CERTIFICATE"
	blockKey	= " KEY"
	blockPub	= " PUBLIC"
	blockPriv	= " PRIVATE"
)


func SaveCertificate(certbytes []byte, filename string) error {

	// Public key
	certOut, err := os.Create(filename)
	if nil != err {
		return err
	}
	
 	pem.Encode(certOut, &pem.Block{Type: blockCert, Bytes: certbytes})	// XXX

	certOut.Close()
	
	return nil
}

func SavePrivateKey(data *rsa.PrivateKey, filename string, pc cakit.PasswordCallback, kdf cakit.KDF) error {
	
 var password string
 var err error
 
	// If callback supplied, call it!
	if nil!=pc {
		
		password,err = pc()
		if nil != err {
			return err
		}
		
	}
		
	// key
	// certOut, err := os.Create(filename)
	keyOut, err := os.OpenFile(filename+".key", 
						os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if nil!=err {
		return err
	}
	defer keyOut.Close()
	
	
	// Setup "block Type"
	blockType := keyType+blockPriv+blockKey
	
	// ...and marshal data as PKCS#1 Private Key
	keyBytes := x509.MarshalPKCS1PrivateKey(data)
	
	var PEMblock *pem.Block
	if nil!=pc && nil!=kdf && ""!=password {
		
		key,_ := kdf(password)
		
		// Encrypted Private KEY block
		PEMblock,err = x509.EncryptPEMBlock(rand.Reader, blockType, keyBytes, key, x509.PEMCipherAES256)
		if nil!=err {
			return err
		}
		
	} else {
		
		// Unencryted Private KEY block
		PEMblock = &pem.Block{
			Type: blockType,
			Bytes: keyBytes,
		}
	}
	
	err = pem.Encode(keyOut, PEMblock)
	if nil!=err {
		return err
	}
	
	return nil
}

func SavePublicKey(data *rsa.PublicKey, filename string) error {
	
	// key
	// certOut, err := os.Create(filename)
	keyOut, err := os.OpenFile(filename+".key", 
						os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if nil != err {
		return err
	}
	
	blockType := keyType+blockPub+blockKey
	
	pembytes, _ := x509.MarshalPKIXPublicKey(data)
	
	pem.Encode(keyOut, &pem.Block{
			Type: blockType, 
			Bytes: pembytes,
		})
	keyOut.Close()
	
	return nil
}

// Result: DER-encoded PKIXPublicKey raw bytes
func LoadPubKey(fname string) ([]byte,string,error) {
	
	pkPEM, err := ioutil.ReadFile(fname)
	if nil != err {
		return nil,"",errors.New("Could not read pubkey: "+err.Error())
	}
	
	PEMblock,_ := pem.Decode(pkPEM)
	if nil==PEMblock {
		return nil,"",errors.New("Unable to parse PEM")
	}
	
	return PEMblock.Bytes,PEMblock.Type,nil
}

func LoadCertificate(fname string) ([]byte,error) {

	certPEM, err := ioutil.ReadFile(fname)
	if nil!=err {
		return nil,err
	}

	PEMblock,_ := pem.Decode(certPEM)
	if blockCert != PEMblock.Type {
		return nil,errors.New("Unexpected PEM block type found: "+PEMblock.Type)
	}
	
	return PEMblock.Bytes,nil
}

func LoadBundle(fname string) ([]x509.Certificate,error) {
	
	bundlePEM, err := ioutil.ReadFile(fname)
	if nil != err {
		return nil,err
	}
	
 var certList = make([]x509.Certificate,0,2)
 var block *pem.Block
 
	for len(bundlePEM) > 0 {
		
		// Decode PEM blocks sequentially
		
		// pem.Decode([]byte) -> (*pem.Block,[]byte)
		block,bundlePEM = pem.Decode(bundlePEM)
		if nil==block {
			break	// no certs left
		}
		if blockCert != block.Type || len(block.Headers) > 0 {
			// Must be encrypted block; Don't know how to handle, so skip it
			continue
		}
		cert,err := x509.ParseCertificate(block.Bytes)
		if nil!=err {
			continue	// invalid cert bytes? Try next
		}
		
		if nil==cert {
			return nil,errors.New("Invalid certificate read??!")
		}
		
		// Fixup certificate (we use custom extensions...)
		if nil!=cert.UnhandledCriticalExtensions {
			cert.UnhandledCriticalExtensions=nil
		}
		
		certList = append(certList,*cert)
		
	} // whend  [[while len(bundlePEM)>0]]

	return certList,nil
}
