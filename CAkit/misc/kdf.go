// kdf.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>
package misc

import (
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
)


const (
	pbkdf2_iter		= 4096
	pbkdf2_outlen	= 256/8		// Suitable for AES256
	
	salt_text	= "843e82b187ed7ad84645dddbcb5f229c58f9c35efe0ef7737cbe1099bd8d0e80"
)


func PBKDF2_AES256(password string) ([]byte, error) {
	
	key := pbkdf2.Key([]byte(password), []byte(salt_text), pbkdf2_iter, pbkdf2_outlen, sha256.New)
	
	return key,nil
}
