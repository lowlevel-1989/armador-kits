// hmac1.c
// Implementation of HMAC_SHA1 as per RFC2104
// 
#include <stdint.h>
#include <string.h>

#include "sha1.h"
#include "hmac.h"


#define PADSIZE 64

inline
int MyHMAC(byte* result,
		 byte* buf, unsigned len,
		 byte* pkey, unsigned kl,
		 SHA1_CTX* ctx)
{
 uint8_t k_ipad[PADSIZE];
 uint8_t k_opad[PADSIZE];

 	memset(k_ipad, 0, PADSIZE);
 	memset(k_opad, 0, PADSIZE);
	memcpy(k_ipad, pkey, PADSIZE-kl);
	memcpy(k_opad, pkey, PADSIZE-kl);
 
	for(register unsigned i=0;i<PADSIZE/4;i++)
	{
		((uint32_t*)k_ipad)[i] ^= 0x36363636;
		((uint32_t*)k_opad)[i] ^= 0x5c5c5c5c;
	}

	// Compute inner = H(Key XOR ipad, buf))
	SHA1Init(ctx);
	SHA1Update(ctx, k_ipad, PADSIZE);
	SHA1Update(ctx, buf, len);
	SHA1Finalize(result, ctx);

	// Compute outer = H(K XOR opad, *inner*)
	SHA1Init(ctx);
	SHA1Update(ctx, k_opad, PADSIZE);
	SHA1Update(ctx, result, HMSZ);
	SHA1Finalize(result, ctx);
	
	return HMSZ;
}


#ifdef TEST

#include <stdio.h>

int main(void)
{
 char key[]   = "\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c\x0c";
 char data[]  = "Test With Truncation";
 byte result[HMSZ];

	SHA1_CTX ctx;

	MyHMAC(result, data, 20, key, 20, &ctx);

	
	for (unsigned i=0; i < HMSZ; i++)
		printf("%02X", result[i]);
	printf("\n");
	
	return 0;
}

#endif
