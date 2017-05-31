// certs.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>
package CAKit


import (
	"crypto/x509"
 	"crypto/x509/pkix"
	"math/big"
	"time"
)



func NewCertificate(subject pkix.Name, serial *big.Int, SKI []byte,
		data []pkix.Extension,
		validFrom time.Time, validUntil time.Time) *x509.Certificate {
	
	// Prepare certificate "template"
	cert := &x509.Certificate{
		SerialNumber:	serial,
		Subject: 		subject,
		NotBefore:		validFrom,
		NotAfter:		validUntil,
		ExtKeyUsage: 	nil,	//[]x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage: 		x509.KeyUsageDigitalSignature,
		ExtraExtensions:data,		// XXX: marshaled to cert; will read as "Extensions"
	}

	return cert
}
