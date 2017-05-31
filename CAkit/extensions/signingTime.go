// signingTime.go
// (C) 2017 Armador Technologies
// Author: Jose Luis Tallon <jltallon@armador.xyz>
package extensions

import (
// 	"encoding/pem"
// 	"time"
)


/* FROM RFC5652 sec11.3
 *  id-signingTime OBJECT IDENTIFIER ::= { iso(1) member-body(2)
 *         us(840) rsadsi(113549) pkcs(1) pkcs9(9) 5 }
 *
 *  SigningTime ::= CHOICE {
 *       utcTime UTCTime,
 *       generalizedTime GeneralizedTime }
 *
 *
 * Dates between 1 January 1950 and 31 December 2049 (inclusive) MUST be
 * encoded as UTCTime.  Any dates with year values before 1950 or after
 * 2049 MUST be encoded as GeneralizedTime.
 */