#ifndef _HMAC_H__INCLUDED_
#define _HMAC_H__INCLUDED_


#define MyHMAC _mK
#define HMSZ	20

#include "sha1.h"


int MyHMAC(byte* result,
		 byte* buf, unsigned len,
		 byte* pkey, unsigned kl,
		 SHA1_CTX* ctx);


#endif /* _HMAC_H__INCLUDED_ */
