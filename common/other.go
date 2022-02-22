package common

import "math/big"

const (
	Zero = "0"
)

func BitsUnzip(bits uint32) *big.Int {
	mantissa := bits & 0xffffff00
	mantissa >>= 8
	e := uint(bits & 0xff)
	c := uint(3)
	var bn *big.Int
	if e <= c {
		shift := 8 * (c - e)
		mantissa >>= shift
		bn = big.NewInt(int64(mantissa))
	} else {
		bn = big.NewInt(int64(mantissa))
		shift := 8 * (e - c)
		bn.Lsh(bn, shift)
	}
	return bn
}
