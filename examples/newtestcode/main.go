package main

import (
	"fmt"
	"mi/common/crc16"
	"strconv"

	"github.com/shopspring/decimal"
)

func main() {
	// header
	header := []byte("0xAAF2")
	// 卡号
	code := "132643667334121621"
	auto_code := []byte{}
	for _, s := range code {
		sbyte := []byte(strconv.FormatInt(int64(rune(s)), 16))
		auto_code = append(auto_code, sbyte...)
	}
	// 支付金额
	amount := decimal.NewFromFloat(10.45).Mul(decimal.NewFromInt(100)).IntPart()

	payamount := []byte(strconv.FormatInt(amount, 16))
	body := append(payamount, auto_code...)

	// 内容长度
	bodyLength := []byte(strconv.FormatInt(int64(len(body)), 16))

	result := append(header, bodyLength...)
	result = append(result, body...)
	crc16 := []byte(crc16.CRCS(result))
	result = append(result, crc16...)
	fmt.Printf("%s\n", result)
}
