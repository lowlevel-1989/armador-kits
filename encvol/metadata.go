// metadata.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>
package encvol

/*
#include <stdlib.h>
#include <unistd.h>
#include <stdint.h>
#include <string.h>
#include <fcntl.h>


#define EVMETA_SIZE	1048576UL
#define EVMETA_BLK	256
#define EVMETA_PAD	((EVMETA_BLK/sizeof(uint32_t))-16*sizeof(uint32_t))
#define EVMETA_KEY	0x1A
#define CKSUM_LEN	20

struct evmeta {
	uint32_t	_pad[EVMETA_PAD];
	char		tagName[16];
	uint8_t		c;
	uint8_t		x;
	uint16_t	m;
	uint32_t	lut;		// Last Used Time
	uint32_t	lu_hi;
	uint32_t	ct;			// Created Time
	uint32_t	ct_hi;
	uint32_t	sz_lo;
	uint32_t	sz_hi;
	uint8_t		k[CKSUM_LEN];
} __attribute__((packed));

typedef struct evmeta evmeta;



inline size_t	_evmetaSize() {
	return sizeof(struct evmeta);
}
size_t	metaSize() { return EVMETA_SIZE; }

inline
const char* _x(char *x, uint64_t len, uint8_t z)
{
	for(register unsigned i=0; i<len; ++i)
		x[i]^=(char)(z+i*3);
}

int readRand(void* buf, size_t len)
{
	int fd = open("/dev/urandom", O_RDONLY);
	int ret= read(fd,buf,len);
	close(fd);
	return ret;
}

int readMeta(evmeta* pm, int fd)
{
	int q = read(fd, pm,sizeof(evmeta));
	uint8_t c = ~pm->c;	
	_x((char*)&pm->m, sizeof(pm->m), ~c);
	_x((char*)pm->tagName, sizeof(pm->tagName), c);
	_x((char*)&pm->lut, 2*sizeof(pm->lut), c);
	_x((char*)&pm->ct, 2*sizeof(pm->ct), c);
	
	_x((char*)&pm->sz_lo, sizeof(uint64_t), 0x55);
	return q;
}

int writeMeta(evmeta* pm, int fd, uint8_t c)
{
	pm->c = ~c;
	_x((char*)&pm->m, sizeof(uint16_t), ~c);
	_x((char*)pm->tagName, sizeof(pm->tagName), c);
	pm->lu_hi = 0;			// should be an assert ...
	_x((char*)&pm->lut, 2*sizeof(pm->lut), c);
	pm->ct_hi=0;
	_x((char*)&pm->ct, 2*sizeof(pm->ct), c);
	
	_x((char*)&pm->sz_lo, sizeof(uint64_t), 0x55);
	
	return write(fd, pm,sizeof(evmeta));
}


inline void _prep(evmeta* m)
{
	memset(m->k, 0xFF,CKSUM_LEN);
}
inline void _saveCK(evmeta* m, uint8_t* save)
{
 	memcpy(save, m->k,CKSUM_LEN);
 	memset(m->k, 0xFF,CKSUM_LEN);
}
// inline void _storeCK(evmeta* m, uint8_t* buf)
// {
// 	memcpy(m->k, buf, CKSUM_LEN);
// }
*/
import "C"

import (
	"errors"
	"os"
	"time"
	"unsafe"
	
	"armador/secmod"
// 	"fmt"
)

const (
	EVMETA_MAGIC = 0xA53C
)

type ContainerMeta struct {
	TagName		string
	Size 		uint64
	Created		time.Time
	LastUsed	time.Time
}

func ContainerMetaSize() uint64 {
	return uint64(C.metaSize())
}


func readContainerMeta(f *os.File, key []byte, cm *ContainerMeta) error {
	
	off := int64(C.metaSize() - C._evmetaSize())
	
	p,err := f.Seek(off,0);	// seek from the beginning of the file
	if p!=off {
		return errors.New("Malformed/corrupt encvol")
	}
	if nil!=err {
		return err
	}
	
	var evm C.struct_evmeta
	r := C.readMeta(&evm, C.int(f.Fd()))
	// Check sanity / detect corruption or tampering
	if r != C.int(C._evmetaSize()) {
		return errors.New("Could not read metablock")
	}
	if EVMETA_MAGIC != evm.m {
		return errors.New("Corrupted metablock found")
	}
	
	return getMeta(&evm, cm)
}

func getMeta(evm *C.struct_evmeta, cm *ContainerMeta) error {
	
	evm.c = C.uint8_t(0);
	cm.TagName = C.GoString(&evm.tagName[0])
	cm.LastUsed = time.Unix(int64(evm.lut),0)
	cm.Created = time.Unix(int64(evm.ct),0)
	
	cm.Size = uint64(evm.sz_hi<<32 | evm.sz_lo)
	
	return nil
}

// Try & guard against "known plaintext" attacks;
// Though there is no "ripple effect" for modified bits to 
func writeContainerMeta(f *os.File, key []byte, cm *ContainerMeta) error {
	
 var evm C.struct_evmeta
 
	ms := C._evmetaSize()

	// Initialize buffer with random data...
	C.readRand(unsafe.Pointer(&evm),ms)
	
	// ...and then get the job done :)
	evm.m = EVMETA_MAGIC
	evm.x = C.uint8_t(ms)
	
	off := int64(C.metaSize() - ms)
	f.Seek(off,0);	// seek from the beginning of the file
	
	C.memset(unsafe.Pointer(&evm.tagName[0]), 0xFF, 16)
	var psz *C.char = C.CString(cm.TagName)
	C.strncpy(&evm.tagName[0], psz, 16)
	C.free(unsafe.Pointer(psz))
	
	// Assume below 2k38 for now;
	// Need to figure out using "C.uint64_t" from Golang
	evm.lut = C.uint32_t(cm.LastUsed.Unix())
	evm.ct = C.uint32_t(cm.Created.Unix())
	evm.lu_hi = 0;		// XXX: fixme
	evm.ct_hi = 0;
	
	evm.sz_lo = C.uint32_t(cm.Size)
	evm.sz_hi = C.uint32_t(cm.Size>>32)

	////////////////////////////////////////////////////////
	C._prep(&evm)
	secmod.BufMAC(unsafe.Pointer(&evm), uint32(ms), key, unsafe.Pointer(&evm.k))
	////////////////////////////////////////////////////////
	
	_, err := C.writeMeta(&evm, C.int(f.Fd()), C.EVMETA_KEY)
	return err
}

func checkHeader(pevm *C.struct_evmeta, key []byte) bool {

 var ck [C.CKSUM_LEN]byte
 	C._saveCK(pevm, (*C.uint8_t)(&ck[0]))
	
	return secmod.BufCheck(unsafe.Pointer(pevm),
					uint32(C._evmetaSize()), 
					key,
					unsafe.Pointer(&ck[0]))
}
