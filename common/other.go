package common

import (
	"math/big"
	"time"
)

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

func BigByZip(target *big.Int) uint32 {
	if target.Sign() == 0 {
		return 0
	}
	c := uint(3)
	e := uint(len(target.Bytes()))
	var mantissa uint
	if e <= c {
		mantissa = uint(target.Bits()[0])
		shift := 8 * (c - e)
		mantissa <<= shift
	} else {
		shift := 8 * (e - c)
		mantissaNum := target.Rsh(target, shift)
		mantissa = uint(mantissaNum.Bits()[0])
	}
	mantissa <<= 8
	mantissa = mantissa & 0xffffffff
	return uint32(mantissa | e)
}

func GetBeforeTime(_day int) int64 {
	// now := time.Now()
	timeZone := time.FixedZone("CST", 8*3600) // 东八区
	nowTime := time.Now().In(timeZone)
	beforeTime := nowTime.AddDate(0, 0, _day)
	beforeTimes := time.Date(beforeTime.Year(), beforeTime.Month(), beforeTime.Day()+1, 0, 0, 0, 0, beforeTime.Location()).Unix()
	return beforeTimes
}
