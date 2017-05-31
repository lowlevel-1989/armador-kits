package memblock

import (
    "errors"
//    "fmt"
//    "os"
    "testing"
)


func TestAlloc(t *testing.T) {

    ps := PageSize()
    mb := MemBlock_Alloc(ps*2)
//    fmt.Fprintln(os.Stderr,"ps=%u; p=%v", ps,mb.Data())

    if mb.IsNull() {
        t.Fatal(errors.New("Eh!"))
    }

}
