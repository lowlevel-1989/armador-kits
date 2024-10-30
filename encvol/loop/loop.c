/*
 * REF: http://man7.org/linux/man-pages/man4/loop.4.html
 */

#include "loop.h"


static inline int _devfd(char* buf, size_t sz, u_int8_t devnr) {
	snprintf(buf, sz, loopbase, devnr); buf[sz-1]='\0';
	return open(buf, O_RDWR|O_CLOEXEC);
}


int setupNodeX(u_int8_t devnr, char* backing_fn, size_t off, size_t sizelimit, u_int8_t* pbKey, unsigned keyLen) {
	if( !pbKey || keyLen<16 || keyLen>32 ){
		return -1;
	}

	char nbuf[16];
	int lfd, backing_fd;

	if( (lfd=_devfd(nbuf,sizeof(nbuf), devnr)) < 0){
		return -1;
	}

	backing_fd = open(backing_fn, O_RDWR | O_CLOEXEC);

	// Setup backing file
	int err;
	if( (err=ioctl(lfd, LOOP_SET_FD, backing_fd)) < 0 ){
		close(lfd);
		return -1;
	}

	struct loop_info64 li64;
	memset(&li64,0,sizeof(struct loop_info64));
	li64.lo_offset = off;
	li64.lo_sizelimit = sizelimit; // bytes, 0 == max available

  struct utsname system_info;

  // Obtener la información del sistema
  if (uname(&system_info) != 0) {
		close(lfd);
		return -1;
  }

  // Verificar si la versión del kernel es menor o igual a 3
  // int major_version = atoi(system_info.release);
  //if (major_version <= 3) {
	//  li64.lo_encrypt_type=1;
  //}

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


int detachNode(u_int8_t devnr) {
	char buf[16];
	int lfd;

	if( (lfd=_devfd(buf,sizeof(buf), devnr)) < 0) {
		return -1;
	}

	int err;
	if( (err=ioctl(lfd, LOOP_CLR_FD, 0)) < 0 ) {
		(void)close(lfd);
		return -1;
	}

	return close(lfd);
}

int openCN() {
	return open("/dev/loop-control", O_RDWR|O_CLOEXEC);
}

int closeCN(int fd) {
	return close(fd);
}

int loopGetFree(int cfd) {
	return ioctl(cfd, LOOP_CTL_GET_FREE);
}
