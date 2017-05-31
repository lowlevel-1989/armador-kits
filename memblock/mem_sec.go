// mem_sec.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>

package memblock


// As per CERT C Secure Coding rule MEM06-C
// Implementation cfr. Open-Std WG 14, Document N1381, improved by the author
// Input also from http://www.stanford.edu/~blp/papers/shredding.pdf 
// A.K.A. https://benpfaff.org/papers/shredding.html/
/// #cgo CFLAGS: -O1 -fno-unroll-loops

/*
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
void *sec_memset(void* v, size_t n, unsigned char c) {
	volatile unsigned char* p = v;
	volatile unsigned char* z = (unsigned char*)p + n;
	while(n--)*p++ = c ^ (~n&0xFF);
	p=v;
	while(p<z)*p++ = ~c;
	while(z<(unsigned char*)v)*z++ = 0;
	return v;
}*/
import "C"


func (mb *MemBlock) Wipe() {

    if mb.IsNull() { return; }
    C.sec_memset(mb.ptr, C.size_t(mb.blen), 0x55)
}


func (mb *MemBlock) Zero() {
	if mb.IsNull() { return; }
	C.memset(mb.ptr,0,C.size_t(mb.blen))
}

func WipeMem(mb *MemBlock, val uint8) {
    C.sec_memset(mb.ptr,C.size_t(mb.blen),C.uchar(val))
}
