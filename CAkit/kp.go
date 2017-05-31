// kp.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>
package CAKit

import (
	"crypto/rsa"
	"crypto/rand"
)

func CreateKeyPair(bitLength uint) (prv *rsa.PrivateKey, pub *rsa.PublicKey) {
	
	prv, _ = rsa.GenerateKey(rand.Reader, int(bitLength))
	pub = &prv.PublicKey
	return
}
