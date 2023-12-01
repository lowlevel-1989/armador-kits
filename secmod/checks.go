// checks.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>

package secmod


/*
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdint.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <string.h>
#include <errno.h>

// including the implementation suffices -- will include the headers too
#include "lib/sha1.c"
#include "lib/hmac.c"

#include "lib/_strx.inl"

#define HASH_SIZE	HMSZ
#define KEY_SIZE	(256/8)

int 
checkFile(const char* name, const uint8_t* hash, size_t fsize)
{
	int fd = 0;

        // debug
        // printf("HIT   [%s]\n", name);
        // char command[100];
        // sprintf(command, "ls -l %s", name);
        // system(command);

	if( (fd=open(name, O_RDONLY|O_CLOEXEC,444)) < 0 )
        {
         // printf("Error: %s\n", strerror(errno));
	 return 1;
        }

	SHA1_CTX ctx;
	SHA1Init(&ctx);

	struct stat sb;
	if( fstat(fd,&sb) < 0 )
	{
		close(fd);
		return 2;
	}

	// check binary hash before execvp
	uint8_t* ph = alloca(HASH_SIZE);
	const char* pd;

	if( !(pd=mmap(NULL,sb.st_size,PROT_READ,MAP_PRIVATE,fd,0)) )
	{
		close(fd);
		return 3;
	}

	SHA1Update(&ctx,(const unsigned char*)pd,sb.st_size);
	SHA1Finalize(ph,&ctx);

	munmap((void*)pd,sb.st_size);

	close(fd);


	_x(ph,HASH_SIZE,0xA5);

        // TODO: habilitar la validación del tamaño de archivo
	// if( (sb.st_size ^ fsize) || memcmp(ph,hash,HASH_SIZE) )

        // se valida el hash
	if( memcmp(ph,hash,HASH_SIZE) )
	return 5;

	return 0;
}

int hsignBlock(byte* res, byte *buf, unsigned len, byte* pkey)
{
	SHA1_CTX ctx;	// func local.....

	MyHMAC(res, buf,len, pkey,KEY_SIZE, &ctx);

	return 0;
}

int checkBlock(byte* px, byte* buf, unsigned len, byte* pkey)
{
	SHA1_CTX ctx;	// func local.....
	uint8_t res[HASH_SIZE];

	MyHMAC(res, buf,len, pkey,KEY_SIZE, &ctx);

	memset(&ctx, 0, sizeof(SHA1_CTX));

	return memcmp(res, px, HASH_SIZE);
}
*/
import "C"

import "unsafe"

func CheckNextFile(name string, hash []byte, filesize uint64) bool {

	ph := C.CBytes(hash)
	defer C.free(ph)

	var r C.int
	var err error

	if r,err = C.checkFile(C.CString(name), (*C.uint8_t)(ph), C.size_t(filesize)); nil!=err {
		return false
	}

	return ( 0 == r )
}

// Safe wrapper
func BufferMAC(pBuf unsafe.Pointer, buflen uint32, k []byte, x []byte) bool {

	if nil == pBuf || 32 != len(k) || 20 < cap(x) {
		return false
	}
	pk := unsafe.Pointer(&k[0])
	px := unsafe.Pointer(&x[0])
	C.hsignBlock((*C.uint8_t)(px),
	(*C.uint8_t)(pBuf), C.uint(buflen),
	(*C.uint8_t)(pk))

	return true
}

// Safe wrapper
func BufferCheck(pBuf unsafe.Pointer, buflen uint32, k []byte, x []byte) bool {

	if nil==pBuf || 32!=len(k) || 20<cap(x) {
		return false
	}
	pk := unsafe.Pointer(&k[0])
	px := unsafe.Pointer(&x[0])
	r,_ := C.checkBlock( (*C.uint8_t)(px),
	(*C.uint8_t)(pBuf), C.uint(buflen),
	(*C.uint8_t)(pk))

	return 0 == r
}


//////////////////////////////////////////////////////////////////////////////
// Unsafe, "raw" funcs from here .... :O

func BufMAC(p unsafe.Pointer, len uint32, k []byte, px unsafe.Pointer) {

	pk := unsafe.Pointer(&k[0])
	C.hsignBlock((*C.uint8_t)(px),
	(*C.uint8_t)(p), C.uint(len),
	(*C.uint8_t)(pk))
}

func BufCheck(pBuf unsafe.Pointer, len uint32, k []byte, px unsafe.Pointer) bool {

	pk := unsafe.Pointer(&k[0])
	r,_ := C.checkBlock( (*C.uint8_t)(px),
	(*C.uint8_t)(pBuf), C.uint(len),
	(*C.uint8_t)(pk))

	return 0 == r
}
