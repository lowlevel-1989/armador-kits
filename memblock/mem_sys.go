// mem_sys.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>

package memblock

// #include <unistd.h>
import "C"

func PageSize() uint32 {

    return uint32(C.sysconf(C._SC_PAGESIZE))

}
