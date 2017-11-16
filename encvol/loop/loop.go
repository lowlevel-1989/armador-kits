// loop.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>

package loop

/*
#include <linux/loop.h>
#include <sys/ioctl.h>
#include <sys/stat.h>
#include <stdlib.h>
#include <unistd.h>
#include <stdint.h>
#include <string.h>
#include <stdio.h>
#include <errno.h>
#include <fcntl.h>

// MUSL does not include <sys/sysmacros.h> automatically (for makedev)....
// ...so supply our own
#ifndef makedev

  #define makedev(x,y) ( \
	(((x)&0xfffff000ULL) << 32) | \
	(((x)&0x00000fffULL) << 8) | \
	(((y)&0xffffff00ULL) << 12) | \
	(((y)&0x000000ffULL)) )

#endif

static const char loopbase[] = "/dev/loop%hu";

int createNode(uint8_t devnr)
{
 char buf[16];
 
	snprintf(buf, sizeof(buf), loopbase, devnr);
// 	buf[sizeof(buf)-1]='\0';
	return mknod(buf, S_IFBLK | 0660, makedev(7,devnr));
}

int openCN()
{	
	return open("/dev/loop-control", O_RDWR|O_CLOEXEC);	// | O_NONBLOCK
}

int loopGetFree(int cfd)
{
	return ioctl(cfd, LOOP_CTL_GET_FREE);
}

int loopAdd(int cfd, uint8_t devnr)
{
	long d = devnr;
	return ioctl(cfd, LOOP_CTL_ADD, d);
}

int loopRemove(int cfd, uint8_t devnr)
{
	long d = devnr;
	return ioctl(cfd, LOOP_CTL_REMOVE, d);
}

int setDIO(int fd, int use_dio)
{
#ifdef LOOP_SET_DIRECT_IO
	// Kernels prior to v4.4 don't support this ioctl (!)
	return ioctl(fd, LOOP_SET_DIRECT_IO, use_dio);
#else
	return 0
#endif
}

int closeCN(int fd)
{
	return close(fd);
}


static inline int
_devfd(char* buf, size_t sz, uint8_t devnr)
{
	snprintf(buf, sz, loopbase, devnr); buf[sz-1]='\0';
	return open(buf, O_RDWR|O_CLOEXEC);
}

int setupNode(uint8_t devnr, int backing_fd, size_t off)
{
 char buf[16];
 int lfd;
	if( (lfd=_devfd(buf,sizeof(buf), devnr)) < 0)
		return -1;
	
	// Setup backing file
 int err;
	if( (err=ioctl(lfd, LOOP_SET_FD, backing_fd)) < 0 )
	{
		close(lfd);
		return -1;
	}
	
 struct loop_info64 li64;
	memset(&li64,0,sizeof(struct loop_info64));
	li64.lo_offset = off;
	
// 	if( 0 != temp )
// 	{
// 		li64.lo_flags = LO_FLAGS_AUTOCLEAR;
// 	}
	
	// Do set "status" (offset+flags)
	if( (err=ioctl(lfd, LOOP_SET_STATUS64, &li64)) < 0 )
	{
		// Attempt cleanup
		(void)ioctl(lfd, LOOP_CLR_FD, 0);
		
		close(lfd);
		return -1;
	}
		
// 	if( 0 != temp )
// 		close(lfd);
	
	return lfd;
}

#ifdef LOOP_SETKEY
int setKey(int fd, uint8_t* pbKey, unsigned keyLen)
{
 struct loop_info64 li64;

	if( !pbKey || 16<keyLen || keyLen>32 )
	{
		errno=EINVAL;
		return -1
	}
	
	if( ioctl(fd, LOOP_GET_STATUS64, &li64) < 0 )
		return -1;

	li64.lo_encrypt_type=1;
	
	unsigned kl = (keyLen>32)? keyLen % 32 : keyLen;
	li64.lo_encrypt_key_size=kl;
	
	memset(li64.lo_encrypt_key, 0,sizeof(li64.lo_encrypt_key));
	memcpy(li64.lo_encrypt_key, pbKey,kl);

	if( ioctl(loopfd, LOOP_SET_STATUS64, &li64) < 0 )
		return -1;

	return 0;
}
#endif

int setupNodeX(uint8_t devnr, int backing_fd, size_t off, uint8_t* pbKey, unsigned keyLen)
{
	if( !pbKey || keyLen<16 || keyLen>32 )
	{
		errno=EINVAL;
		return -1;
	}

 char nbuf[16];
 int lfd;

	if( (lfd=_devfd(nbuf,sizeof(nbuf), devnr)) < 0)
		return -1;

	// Setup backing file
	int err;
	if( (err=ioctl(lfd, LOOP_SET_FD, backing_fd)) < 0 )
	{
		close(lfd);
		return -1;
	}

 struct loop_info64 li64;
	memset(&li64,0,sizeof(struct loop_info64));
	li64.lo_offset = off;
	
	li64.lo_encrypt_type=1;
	
	unsigned kl = (keyLen>32)? keyLen % 32 : keyLen;
	li64.lo_encrypt_key_size=kl;
	
	memset(li64.lo_encrypt_key, 0,sizeof(li64.lo_encrypt_key));
	memcpy(li64.lo_encrypt_key, pbKey,kl);

	// Do set "status" (offset+flags+key)
	if( (err=ioctl(lfd, LOOP_SET_STATUS64, &li64)) < 0 )
	{
		// Attempt cleanup
		(void)ioctl(lfd, LOOP_CLR_FD, 0);

		close(lfd);
		return -1;
	}

	return lfd;
}

int detachNode(uint8_t devnr)
{
 char buf[16];
 int lfd;
	if( (lfd=_devfd(buf,sizeof(buf), devnr)) < 0)
		return -1;
 
  int err;
	if( (err=ioctl(lfd, LOOP_CLR_FD, 0)) < 0 )
	{
		(void)close(lfd);
		return -1;
	}
	
	return close(lfd);
}

int LoopAlloc(int cfd, uint8_t min)
{
 long devnr, ret;
 
	devnr = ioctl(cfd, LOOP_CTL_GET_FREE);
	if( devnr < 0 )
		return -1;
	
	if( devnr < min )
	{
		// Release the loopdev, since it doesn't fulfill reqs
		ioctl(cfd, LOOP_CTL_REMOVE, devnr);
		
		devnr=min;
		do {
			errno=0;
			devnr++;
			ret = ioctl(cfd, LOOP_CTL_ADD, devnr);
		} while( ret<0 && EEXIST==errno );
		
		devnr=ret;
	}
	
	return (devnr<0) ? -1 : (devnr & 0x000000FF);
}

// static 
// char* loop_devName(uint8_t* buf, unsigned len, uint8_t devnr)
// {
// 	snprintf(buf, len, loopbase, devnr);
// 	return buf;
// }
*/
import "C"


import (
	"os"
)

var (
	cfd	int	= -1
)

////////////////////////////////////////////////////////////////////////////////

// func init() {
// 
// 	r,_ := C.openCN()
// 	cfd = int(r)
// }

////////////////////////////////////////////////////////////////////////////////

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

func Allocate(min uint8) (uint8,error) {
	
// 	cfd,err := C.openCN()
// 	if cfd < 0 || nil!=err {
// 		return 0xFF,err
// 	}
// 	defer C.closeCN(cfd)
	dn,err := C.LoopAlloc(C.int(cfd), C.uint8_t(min))	
	if dn < 0 || nil!=err {
		return 0xFF,err
	}
	return uint8(dn),nil
}

func CreateNode(devnr uint8) error {
	_,err := C.createNode(C.uint8_t(devnr))
	return err
}

func Assign(devnr uint8) error {
	
	_,err := C.loopAdd(C.int(cfd),C.uint8_t(devnr))
	return err
}

func Release(devnr uint8) error {
	
	_,err := C.loopRemove(C.int(cfd),C.uint8_t(devnr))
	return err
}

func Setup(devnr uint8, bf *os.File, loff uint64) error {
	
	_,err := C.setupNode(C.uint8_t(devnr), C.int(bf.Fd()), C.size_t(loff))
	if nil!=err {
		return err
	}
	
	return nil
}

func SetupX(devnr uint8, bf *os.File, loff uint64, key []byte) error {

	ptr := C.CBytes(key)
	r,err := C.setupNodeX(C.uint8_t(devnr), C.int(bf.Fd()), C.size_t(loff),
					(*C.uint8_t)(&key[0]), C.uint(len(key)))
 	C.free(ptr)
	if r<0 || nil != err {
		return err
	}

	return nil
}

func Attach(devnr uint8, bf *os.File, loff uint64) (int,error) {
	
	r,err := C.setupNode(C.uint8_t(devnr), C.int(bf.Fd()), C.size_t(loff))
	if nil!=err {
		return -1,err
	}
	
	return int(r),nil
}

func Detach(devnr uint8) error {

	_,err := C.detachNode(C.uint8_t(devnr))
	return err;
}


func RawAlloc() (uint8,error) {
	
// 	cfd,err := C.openCN()
// 	if cfd < 0 || nil!=err {
// 		return 0xFF,err
// 	}
// 	defer C.closeCN(cfd)
	dn, err := C.loopGetFree(C.int(cfd))
	if dn < 0 {
		return 0xFF,err
	}

	return uint8(dn),nil
}
