// closures.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>
package CAKit


import (
// 	"fmt"
)

type PasswordCallback func() (string,error)

type KDF func(password string) ([]byte,error)
