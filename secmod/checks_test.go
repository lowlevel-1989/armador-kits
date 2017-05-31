package secmod

import (
	"testing"
	"fmt"
	"unsafe"
)

const (
	
	test_key	= "\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0a\x0a\x0a\x0a\x0a\x0a\x0a\x0a\x0a\x0a\x0b\x0b\x0b\x0b\x0b\x0b\x0b\x0b\x0b\x0b";
	test_data	= "En un lugar de la mancha de cuyo nombre no quiero acordarme vivía un hidalgo caballero de nombre Don Quijote de la Mancha"
)

var (
	key [32]byte
)

func ExampleMAC() {
	
	var res [20]byte
	
	data := []byte(test_data)
	copy(key[:],[]byte(test_key))
	
	BufferMAC(unsafe.Pointer(&data[0]), uint32(len(data)), key[:], res[:])
	
	fmt.Printf("%x\n",res)
	// Output: 8489412a027ffd98d3378e9f8431da4debf9545e
}


func TestMain(m *testing.M) {

	// No init yet
	
}



func TestBlockMAC(t *testing.T) {

	var res [20]byte
	
	data := []byte(test_data)
	copy(key[:],[]byte(test_key))
	
	BufferMAC(unsafe.Pointer(&data[0]), uint32(len(data)), key[:], res[:])

	fmt.Printf("%x\n",res)
}



func TestCheckMAC(t *testing.T) {
	

	var xs string = "\x84\x89A*\x02\u007f\xfd\x98\xd37\x8e\x9f\x841\xdaM\xeb\xf9T^"
	
	data := []byte(test_data)
	copy(key[:],[]byte(test_key))	// discards final '\0'
	
	var x [20]byte
	copy(x[:],xs)
	
	r := BufferCheck(unsafe.Pointer(&data[0]), uint32(len(data)), key[:], x[:])
	
	if !r {
		t.Fail()
	}
}
