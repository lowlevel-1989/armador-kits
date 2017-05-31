// serial.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>
package util

import (
	"crypto/rand"
	"errors"
	"math/big"
)

func GenSerial() (*big.Int,error) {
	
	upperLim := new(big.Int).Lsh(big.NewInt(1), 128)
	result, err := rand.Int(rand.Reader, upperLim)
	if nil!=err {
		return nil,errors.New("Failed to generate serial number:"+err.Error())
	}
	return result,nil
}
