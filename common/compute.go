package common

import (
	"math/big"

	"github.com/shopspring/decimal"
)

func CalcGasFee(x, y float64) *big.Float {
	result := new(big.Float)
	xs := new(big.Float).SetFloat64(x)
	ys := new(big.Float).SetFloat64(y)
	result.Mul(xs, ys)
	return result
}

func Div(x, y int64) decimal.Decimal {
	xs := decimal.NewFromInt(x)
	ys := decimal.NewFromInt(y)
	return xs.Div(ys).Round(2)
}
