// output.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>
package CAKit

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"strings"
)

var (
	pkAlgoNames = map[x509.PublicKeyAlgorithm]string {
		x509.UnknownPublicKeyAlgorithm: "<unknown>",
		x509.RSA:	"RSA",
		x509.DSA:	"DSA",
		x509.ECDSA:	"ECDSA",
	}
)

/*
Raw                     []byte // Complete ASN.1 DER content (certificate, signature algorithm and signature).
RawTBSCertificate       []byte // Certificate part of raw ASN.1 DER content.
RawSubjectPublicKeyInfo []byte // DER encoded SubjectPublicKeyInfo.
RawSubject              []byte // DER encoded Subject
RawIssuer               []byte // DER encoded Issuer

Signature          []byte
SignatureAlgorithm SignatureAlgorithm

PublicKeyAlgorithm PublicKeyAlgorithm
PublicKey          interface{}

Version             int
SerialNumber        *big.Int
Issuer              pkix.Name
Subject             pkix.Name
NotBefore, NotAfter time.Time // Validity bounds.
KeyUsage            KeyUsage

// Extensions contains raw X.509 extensions. When parsing certificates,
// this can be used to extract non-critical extensions that are not
// parsed by this package. When marshaling certificates, the Extensions
// field is ignored, see ExtraExtensions.
Extensions []pkix.Extension

ExtKeyUsage        []ExtKeyUsage           // Sequence of extended key usages.
UnknownExtKeyUsage []asn1.ObjectIdentifier // Encountered extended key usages unknown to this package.

// BasicConstraintsValid indicates whether IsCA, MaxPathLen,
// and MaxPathLenZero are valid.
BasicConstraintsValid bool
IsCA                  bool

SubjectKeyId   []byte
AuthorityKeyId []byte

// RFC 5280, 4.2.2.1 (Authority Information Access)
OCSPServer            []string
IssuingCertificateURL []string

// Subject Alternate Name values
DNSNames       []string
EmailAddresses []string
IPAddresses    []net.IP

// Name constraints
PermittedDNSDomainsCritical bool // if true then the name constraints are marked critical.
PermittedDNSDomains         []string
ExcludedDNSDomains          []string

// CRL Distribution Points
CRLDistributionPoints []string

PolicyIdentifiers []asn1.ObjectIdentifier
}
*/



func DumpCertificate(cert *x509.Certificate) string {
	
 var buf bytes.Buffer
	
	buf.WriteString(fmt.Sprintf("Version: %d\n", cert.Version))
	buf.WriteString(fmt.Sprintf("Serial Number: %x\n", cert.SerialNumber))
	buf.WriteString(fmt.Sprintf("Signature Algorithm: %s\n", cert.SignatureAlgorithm))
	
	buf.WriteString(fmt.Sprintf("Issuer: %s\n", PKIXName_DN(&cert.Issuer)))
	buf.WriteString(fmt.Sprintf("\tAKI: %X\n", cert.AuthorityKeyId))
	buf.WriteString("Validity:\n")
	buf.WriteString(fmt.Sprintf("\tNot Before: %v\n", cert.NotBefore))
	buf.WriteString(fmt.Sprintf("\tNot After:  %v\n", cert.NotAfter))
	
	buf.WriteString(fmt.Sprintf("Subject: %s\n", PKIXName_DN(&cert.Subject)))
	buf.WriteString(fmt.Sprintf("Public Key: "))		// pkAlgoNames[cert.PublicKeyAlgorithm]))
	buf.WriteString(PubKeyToString(cert.PublicKey))
	buf.WriteString(fmt.Sprintf("\n\tSKI: %X\n", cert.SubjectKeyId))
// 	buf.WriteString("\n")
	
	buf.WriteString("Signature:\n")
	buf.WriteString(fmt.Sprintf("\tAlgorithm: %s\n", cert.SignatureAlgorithm))
	buf.WriteString(fmt.Sprintf("\tRaw: %X\n", cert.Signature))
	
	buf.WriteString("Constraints:\n")
	buf.WriteString(fmt.Sprintf("\tIsCA: %v", cert.IsCA))
	
	return buf.String()
}

/*
type Name struct {
	Country, Organization, OrganizationalUnit []string
	Locality, Province                        []string
	StreetAddress, PostalCode                 []string
	SerialNumber, CommonName                  string
	
	Names      []AttributeTypeAndValue
	ExtraNames []AttributeTypeAndValue
}*/
func PKIXName_DN(n *pkix.Name) string {
	
	res := make([]string,0,9)
	mergeStringSlice(&res,n.Country,"C")
	mergeStringSlice(&res,n.Province,"ST")
	mergeStringSlice(&res,n.Locality,"L")
	mergeStringSlice(&res,n.Organization,"O")
	mergeStringSlice(&res,n.OrganizationalUnit,"OU")
	
	res = append(res,fmt.Sprintf("CN = %s", n.CommonName))
	
	// XXX: Ignores "SerialNumber", at least for now
	
	return strings.Join(res,", ")
}

func PKIXName_CN(n *pkix.Name) string {
	
	cn := n.CommonName
	if strings.IndexRune(cn,',') >= 0 {
		cn = `"`+n.CommonName+`"`
	}
	return ("CN = "+cn)
}


func PubKeyToString(pubKey interface{}) string {

	if d,ok := pubKey.(*rsa.PublicKey); ok {
		m:=d.N.Bytes()
		return fmt.Sprintf("RSA (%d bits) {Exponent: %d (0x%X)\nModulus: %X}", len(m)*8, d.E, d.E, m)
	}
	return fmt.Sprintf("%v",pubKey)
}


func mergeStringSlice(result *[]string, data []string, field string) {
	
	if nil==result || nil==data { return; }
	
	buf := make([]string,0,1)
	for _,v := range data {
		if strings.IndexRune(v,',') >= 0 {
			v = `"`+v+`"`
		}
		buf=append(buf, (field+" = "+v))
	}
	dn := strings.Join(buf,", ")
	
	*result = append(*result,dn)
}
