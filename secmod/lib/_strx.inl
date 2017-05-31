#ifndef byte
 #define byte uint8_t
#endif

inline
const char* 
_x(byte *x, size_t len, uint8_t z)
{
	for(register unsigned i=0; i<len; ++i)
	{
		x[i]^=((z+i)&0xFF);
	}
	
	return (const char*)x;
}
