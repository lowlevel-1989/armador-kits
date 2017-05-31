// file_sec.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>

package secmod


/*
#include <stdint.h>

inline
unsigned memtest(uint8_t* restrict buf, unsigned len, uint8_t val)
{
 register unsigned ax;
	ax = (len ^ len);	// => 0 ;)
	while(len--)
		ax |= (*buf++ ^ val);
	return ax;
}


uint8_t* memxor(uint8_t *restrict dst, const uint8_t *restrict src, unsigned n)
{
	for (; n > 0; n--)
		*dst++ ^= *src++;

	return dst;
}

*/
import "C"
