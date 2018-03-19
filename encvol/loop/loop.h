#include <fcntl.h>
#include <linux/loop.h>
#include <sys/types.h>
#include <sys/ioctl.h>
#include <sys/sysmacros.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>

static const char loopbase[] = "/dev/loop%hu";

static inline int _devfd(char* buf, size_t sz, u_int8_t devnr);

int setupNodeX(u_int8_t devnr, char* backing_fn, size_t off, size_t sizelimit, u_int8_t* pbKey, unsigned keyLen);

int detachNode(u_int8_t devnr);

int createNode(u_int8_t devnr);

int openCN();
int closeCN(int fd);
int loopGetFree(int cfd);
