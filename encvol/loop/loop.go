package loop

//#include "loop.h"
//#cgo CFLAGS: -I/usr/include
import "C"

import (
	"path/filepath"
	"os"
	"bufio"
	"strings"
	"strconv"
)

const (
	f_mounts = "/proc/mounts"
)

type MountDev struct {
	DEV string
	PATH string
	FS string
	N uint8
}

var (
	cfd int	= -1 // variable de control
)

func InitAdvanced() error {

	r,err := C.openCN()
	if nil!=err {
		return err
	}
	cfd = int(r)
	return nil
}

func Deinit() (err error) {
	_,err = C.closeCN(C.int(cfd))
	return
}

func RawAlloc() (uint8,error) {
	dn, err := C.loopGetFree(C.int(cfd))
	if dn < 0 {
		return 0xFF,err
	}

	return uint8(dn),nil
}


func SetupX(devnr uint8, backing_fn string, loff uint64, slimit uint64, key []byte) (int, error) {

	ptr := C.CBytes(key)
	r,err := C.setupNodeX(C.u_int8_t(devnr), C.CString(backing_fn), C.size_t(loff),
				C.size_t(slimit), (*C.u_int8_t)(&key[0]), C.uint(len(key)))
	C.free(ptr)
	if r<0 || nil != err {
		return -1, err
	}

	return int(r), nil
}

func Detach(devnr uint8, fddev int) error {

	var err error

	if fddev >= 0 {
		_, err = C.closeCN(C.int(fddev))
		if err != nil {
			return err
		}
	}

	_, err = C.detachNode(C.u_int8_t(devnr))
	return err
}

func GetLoopMount(dir string) ([]MountDev, int) {
	volume_dir, _ := filepath.Abs(dir)

	list_mounts := make([]MountDev, 0, 5)

	f, _ := os.Open(f_mounts)
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, volume_dir) {
			dev_data := strings.Split(line, " ")
			N, _ := strconv.Atoi(string(dev_data[0][len(dev_data[0])-1]))
			dev_loop := MountDev{dev_data[0], dev_data[1], dev_data[2], ((uint8)(N))}
			list_mounts = append(list_mounts, dev_loop)
		}
	}

	return list_mounts, len(list_mounts)

}
