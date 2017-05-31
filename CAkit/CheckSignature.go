// CheckSignature.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>
// Mostly ripped from "crypto/x509/x509.go", which is "(C)2009 The Go Authors."

package CAKit

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"math/big"
)

const (
	msgUnsupportedAlgorithm = "Cannot verify signature: algorithm unimplemented"
)



////////////////////////////////////////////////////////////////////////////////

type ecdsaSignature struct {	
	R, S *big.Int
	
}

func isRSAPSS(algo x509.SignatureAlgorithm) bool {	
	switch algo {
		case x509.SHA256WithRSAPSS, x509.SHA384WithRSAPSS, x509.SHA512WithRSAPSS:
			return true
		default:
			return false
	}
}


func ChkSig(cert *x509.Certificate, publicKey crypto.PublicKey) (err error) {

	if nil==cert {
		return errors.New("Invalid arguments")
	}
	
	var hashType crypto.Hash

	algo := cert.SignatureAlgorithm
	switch algo {
//	// ** OBSOLETE => Unsupported **
// 	case x509.SHA1WithRSA, x509.DSAWithSHA1, x509.ECDSAWithSHA1:
// 		hashType = crypto.SHA1

	case x509.SHA256WithRSA, x509.SHA256WithRSAPSS, x509.DSAWithSHA256, x509.ECDSAWithSHA256:
		hashType = crypto.SHA256

	case x509.SHA384WithRSA, x509.SHA384WithRSAPSS, x509.ECDSAWithSHA384:
		hashType = crypto.SHA384

	case x509.SHA512WithRSA, x509.SHA512WithRSAPSS, x509.ECDSAWithSHA512:
		hashType = crypto.SHA512

	default:
		return errors.New(msgUnsupportedAlgorithm)
	}

	if !hashType.Available() {
		return errors.New(msgUnsupportedAlgorithm)
	}

	signed := cert.RawTBSCertificate
	
	h := hashType.New()
	h.Write(signed)
	digest := h.Sum(nil)

	switch pub := publicKey.(type) {

	case *rsa.PublicKey:
		if isRSAPSS(algo) {
			return rsa.VerifyPSS(pub, hashType, digest, cert.Signature, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash})
		} else {
			return rsa.VerifyPKCS1v15(pub, hashType, digest, cert.Signature)
		}

//	// *** OBSOLETE ***
// 	case *dsa.PublicKey:
// 		dsaSig := new(dsaSignature)
// 		if rest, err := asn1.Unmarshal(signature, dsaSig); err != nil {
// 			return err
// 		} else if len(rest) != 0 {
// 			return errors.New("trailing data after DSA signature")
// 		}
// 
// 		if dsaSig.R.Sign() <= 0 || dsaSig.S.Sign() <= 0 {
// 			return errors.New("DSA signature contained zero or negative values")
// 		}
// 
// 		if !dsa.Verify(pub, digest, dsaSig.R, dsaSig.S) {
// 			return errors.New("DSA verification failure")
// 		}
// 		return

	case *ecdsa.PublicKey:
		ecdsaSig := new(ecdsaSignature)
		if rest, err := asn1.Unmarshal(cert.Signature, ecdsaSig); err != nil {
			return err
		} else if len(rest) != 0 {
			return errors.New("x509: trailing data after ECDSA signature")
		}

		if ecdsaSig.R.Sign() <= 0 || ecdsaSig.S.Sign() <= 0 {
			return errors.New("x509: ECDSA signature contained zero or negative values")
		}

		if !ecdsa.Verify(pub, digest, ecdsaSig.R, ecdsaSig.S) {
			return errors.New("x509: ECDSA verification failure")
		}

		return
	}

	return errors.New(msgUnsupportedAlgorithm)
}
