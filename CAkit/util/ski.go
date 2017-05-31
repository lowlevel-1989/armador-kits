// ski.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>
package util


// XXX: TODO: build database of issuers/issued -> CRL


import (
	"hash"
    "crypto/sha1"
	"crypto/x509"
	"crypto/rsa"
)


/* FROM RFC5280, sec6 -- path building
 * To facilitate certification path construction, this extension MUST appear
 * in all conforming CA certificates, that is, all certificates including the
 * basic constraints extension (Section 4.2.1.9) where the value of CA is TRUE.
 * In conforming CA certificates, the value of the subject key identifier MUST
 * be the value placed in the key identifier field of the authority key
 * identifier extension (Section 4.2.1.1) of certificates issued by the subject
 * of this certificate. Applications are not required to verify that key
 * identifiers match when performing certification path validation.
 */


func SKIfromPub(data []byte) []byte {

    hf := sha1.New()
	hf.Write(data)
    return hf.Sum(nil)
}

func ComputeSKI(pubkey *rsa.PublicKey, hf hash.Hash) ([]byte,error) {
	
 var kb []byte
 
	kb,err := x509.MarshalPKIXPublicKey(pubkey)
	if nil != err {
		return nil,err
	}
	hf.Write(kb)
	return hf.Sum(nil), nil	
}
