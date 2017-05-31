// mem_misc.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>

package memblock


import (
	"reflect"
	"unsafe"
)


// Wrap the MemBlock in a byteSlice, as it is
func (mb *MemBlock) Wrap() []byte {
	
	hdr := &reflect.SliceHeader{uintptr(mb.ptr), int(mb.blen), int(mb.blen)}	
	return *(*[]byte)(unsafe.Pointer(hdr))	
}

// Return a byteSlice backed by the MB; at most "elen" bytes in size
// Equivalent to (([]byte)mb)[:elen]
func (mb *MemBlock) RawBuf(elen uint64) []byte {
	l := elen
	if elen >= mb.blen {
		l = mb.blen
	}
	hdr := &reflect.SliceHeader{uintptr(mb.ptr), int(l), int(mb.blen)}
	return *(*[]byte)(unsafe.Pointer(hdr))	
}

// Return a byteSlice backed by the MB, starting at offset 'off'
// Equivalent to (([]byte)mb)[off:]
func (mb *MemBlock) BufFrom(off uint64) []byte {
	d := off
	if off >= mb.blen {
		d = mb.blen
	}
	hdr := &reflect.SliceHeader{uintptr(mb.ptr)+uintptr(d), int(mb.blen-d), int(mb.blen-d)}
	return *(*[]byte)(unsafe.Pointer(hdr))	
}

// Return a byteSlice backed by the MB, starting at offset 'off'
// Equivalent to (([]byte)mb)[off:elen]
func (mb *MemBlock) Slice(off, elen uint64) []byte {
	
	hdr := &reflect.SliceHeader{uintptr(mb.ptr)+uintptr(off), int(elen), int(mb.blen-off)}
	return *(*[]byte)(unsafe.Pointer(hdr))	
}


func (mb *MemBlock) RawPtr() unsafe.Pointer {
	return mb.ptr
}
