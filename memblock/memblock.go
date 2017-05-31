// memblock.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>

package memblock
 
/* 
#include <unistd.h>
#include <sys/mman.h>
#include <sys/types.h>
#include <stdint.h>

unsigned memtest(uint8_t* src, unsigned len, uint8_t val)
{
 register unsigned ax;
 
    ax = (len ^ len);
    while(len--)
    {
        ax |= (*src++ ^ val);
    }
    return ax;
}
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

type MemBlock struct {
    ptr		unsafe.Pointer
    blen	uint64
}



var (
    err_MbNull error = errors.New("NULL MemBlock !?")
)


func New(bs uint64) (*MemBlock,error) {

	ps := PageSize()
	
    // At least 1 page ; expect mmap to do the right thing w.r.t. alignment
//     elen := C.size_t( (uint32((bs+1)/C._SC_PAGESIZE)+1)*C._SC_PAGESIZE )
	elen := C.size_t( (uint32((bs+1)/uint64(ps))+1) * ps )
    p,e := C.mmap(nil, C.size_t(elen),
            C.PROT_READ|C.PROT_WRITE,
            C.MAP_PRIVATE|C.MAP_ANONYMOUS|C.MAP_POPULATE,
            0,0)
    if nil!=e {
        return &MemBlock{nil,0}, e
    }

    // Linux 3.4+
    r,err := C.madvise(p,C.size_t(elen),C.MADV_DONTDUMP)
    if r<0 || nil!=err {
        _free(p,elen)
        return &MemBlock{nil,0}, err
    }
    // Linux 2.6.19+
    r,err = C.madvise(p,C.size_t(elen),C.MADV_DONTFORK)
    if nil!= err {
        _free(p,elen)
        return &MemBlock{nil,0}, err
    }

    return &MemBlock{p,uint64(elen)}, nil
}

// func Memblock_Alloc2(blen uint64) *MemBlock {
// 
//     p,e := C.mmap(nil, C.size_t(blen),
//             C.PROT_READ|C.PROT_WRITE,
//             C.MAP_PRIVATE|C.MAP_ANONYMOUS|C.MAP_NORESERVE,
//             0,0)
//     if nil!=e {
//         return &MemBlock{nil,0}
//     }
// 
//     return &MemBlock{p,blen}
// }

func (mb *MemBlock) Dispose() error {
    if mb.IsNull() {
        return err_MbNull
    }
 
    WipeMem(mb,0x55)
    WipeMem(mb,0xCC)
    err := _free(mb.ptr, C.size_t(mb.blen))
	
	mb.ptr = nil
	mb.blen = 0
	
	return err
}


func _free(ptr unsafe.Pointer, blen C.size_t) error {
    if nil == ptr { return err_MbNull; }
    if r,e := C.munmap(ptr,blen); r<0 {
        return e
    }
    return nil
}


func (mb *MemBlock) Size() uint64 {

    return mb.blen
}

/*
func (mb *MemBlock) Data() []byte {
	
	return C.GoBytes(mb.ptr,C.int(mb.blen))
}*/

func (mb *MemBlock) String() string {
	
	return fmt.Sprintf("MemBlock[base=%x; blen=%d]", uintptr(mb.ptr), mb.blen)
	
}

func (mb *MemBlock) IsNull() bool {

    return (1 > mb.blen || nil == mb.ptr)
}
