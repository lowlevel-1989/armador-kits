#ifndef _SHA1__H_
#define _SHA1__H_
/*
 * SHA-1 in C by Steve Reid <steve@edmweb.com>
 * Public Domain
 * Modified by Jose Luis Tallon, 2017
 */

#include <stdint.h>

#define SHA1_CTX _hhC
#define SHA1Init	_hhI
#define SHA1Update	_hhU
#define SHA1Finalize	_hhF

typedef struct {
    uint32_t state[5];
    uint32_t count[2];
    unsigned char buffer[64];
} SHA1_CTX;


#ifndef byte
 #define byte uint8_t
#endif


//void SHA1Transform(uint32_t state[5], const unsigned char buffer[64]);
void SHA1Init(SHA1_CTX* context);
void SHA1Update(SHA1_CTX* context, const byte* restrict data, uint32_t len);
void SHA1Finalize(byte digest[20], SHA1_CTX* context);


#endif	// ndef _SHA1__H_
