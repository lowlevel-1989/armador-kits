// customExt.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>
package extensions

import (
	"crypto/x509/pkix"
	"encoding/asn1"
)


func ExtensionFromRaw(oid []int, data []byte) pkix.Extension {

	return pkix.Extension{Id: asn1.ObjectIdentifier(oid), Critical: true, Value: data}

}

func GetExtensionValue(input []pkix.Extension, oid []int) []byte {
	
	for _,x := range input {
		
		if intSliceEquals(oid,x.Id) { return x.Value; }
	}
	return nil
}


func intSliceEquals(x,y []int) bool {
	if nil==x && nil==y { return true; }
	if nil==x || nil==y { return false; }
	
	if len(x) != len(y) { return false; }
	
	for i := range x {
		if x[i] != y[i] { return false; }
	}
	return true
}
