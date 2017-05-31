// memblock.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>

package memblock
 
// #include <unistd.h>
// #include <sys/mman.h>
import "C"

func (mb *MemBlock) Lock() error {

    if mb.IsNull() {
        return err_MbNull
    }
    if t,err := C.mlock(mb.ptr, C.size_t(mb.blen)/*, MLOCK_ONFAULT*/); t < 0 {
        return err
    }
    return nil
}

func (mb *MemBlock) Unlock() error {

    if mb.IsNull() {
        return err_MbNull
    }
    if t,err := C.munlock(mb.ptr, C.size_t(mb.blen)); t<0 {
        return err
    }
    return nil
}

/*
func (mb *MemBlock) InCore() bool {

    if mb.IsNull() {
        return false
    }

    ps := uint64(C._SC_PAGESIZE)
    vlen := uint64((mb.blen+ps-1)/ps)
    // int mincore(void *addr, size_t length, unsigned char *vec);
    //r := C.mincore(mb.ptr, C.size_t(mb.blen), nil)

    return false
}
*/

func (mb *MemBlock) SetNoDump() error {

    if r,e := C.madvise(mb.ptr, C.size_t(mb.blen), C.MADV_DONTDUMP); r<0 {
        return e
    }
    return nil
}
