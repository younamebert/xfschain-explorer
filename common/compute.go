package common

import "math/big"

func CalcGasFee(x, y float64) *big.Float {
	result := new(big.Float)
	xs := new(big.Float).SetFloat64(x)
	ys := new(big.Float).SetFloat64(y)
	result.Mul(xs, ys)
	return result
}
