// password_cb.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>
package util


import (
	"errors"
	"fmt"
	"io"
)

func BasicPasswordCallback(r io.Reader) (string,error) {
	
	var pswd string
	nr,err := fmt.Fscanln(r,&pswd)
	if nil!=err { return "",err; }
	if nr<1 || len(pswd) < 8 {
		return "",errors.New("Truncated password")
	}
	return pswd,nil
}
