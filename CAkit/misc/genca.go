// genca.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>
package misc

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

func RootCACertificate(caname *pkix.Name, validFrom time.Time,validYears uint8) x509.Certificate {

	return x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:	*caname,
		NotBefore:	validFrom,
		NotAfter:	validFrom.AddDate(int(validYears), 0, 0),
		IsCA:		true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:	x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
}


func CACertificate(caname *pkix.Name, serial uint64, validFrom time.Time, validMonths uint8) x509.Certificate {
	
	return x509.Certificate{
		SerialNumber: big.NewInt(int64(serial)),
		Subject:	*caname,
		NotBefore:	validFrom,
		NotAfter:	validFrom.AddDate(0,int(validMonths), 0),
		IsCA:		true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:	x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
}
