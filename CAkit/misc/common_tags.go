package misc

import (
	"encoding/asn1"
//	"crypto/x509/pkix"
	"fmt"
)


type DRMtag uint8

const (
	TagBase 			DRMtag 	= iota
	TagLicenseRuns
)



func DRMtag2ASN1(t uint8, u DRMtag, v uint8) asn1.ObjectIdentifier {

 var firstElem uint8
 
	firstElem = (40*asn1.ClassPrivate + asn1.TagInteger)
 
	return asn1.ObjectIdentifier{int(firstElem), int(t),int(u),int(v)}
}


func DRMtag2String(oid *asn1.ObjectIdentifier) string {
	
	return fmt.Sprintf("")

}
