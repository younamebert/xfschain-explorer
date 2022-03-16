package apis

import (
	"encoding/hex"
	"fmt"

	"github.com/shopspring/decimal"
)

func retuByte(money, b_price float64) []byte {

	//a仓价格
	acs := make([]byte, 0)

	queryPrice := money
	r := decimal.NewFromFloat(100)
	s := r.Mul(decimal.NewFromFloat(queryPrice))

	rs := hex.EncodeToString(s.BigInt().Bytes())
	_ = rs

	if len(s.BigInt().Bytes()) < 2 {
		t := s.BigInt().Bytes()[0]
		fmt.Printf("%T\n", t)
		// ts := []byte{0x00, t}
		acs = append(acs, 0x00, t)
		// asc := hex.EncodeToString(ac)
		// fmt.Println(ts)

	} else {
		t := s.BigInt().Bytes()[0]
		t2 := s.BigInt().Bytes()[1]

		acs = append(acs, t, t2)
	}

	queryPrice_b := b_price
	r2 := decimal.NewFromFloat(100)
	s2 := r2.Mul(decimal.NewFromFloat(queryPrice_b))

	rs2 := hex.EncodeToString(s2.BigInt().Bytes())
	_ = rs2

	if len(s2.BigInt().Bytes()) < 2 {
		t2 := s2.BigInt().Bytes()[0]
		fmt.Printf("%T\n", t2)
		// ts := []byte{0x00, t}
		acs = append(acs, 0x00, t2)
		// asc := hex.EncodeToString(ac)
		// fmt.Println(ts)

	} else {
		t := s2.BigInt().Bytes()[0]
		t3 := s2.BigInt().Bytes()[1]

		acs = append(acs, t, t3)
	}

	// 输出变量的十六进制形式和十进制值

	// c := hex.EncodeToString(acs[1:2])

	fmt.Println(fmt.Sprintf("%#[1]x %#x", acs[:1], acs[1:2]))
	// so := hex.EncodeToString(acs[:1])

	fmt.Println("第一个数:", hex.EncodeToString(acs[:1]))
	fmt.Println("第二个数:", hex.EncodeToString(acs[1:2]))

	fmt.Println("第三个数:", hex.EncodeToString(acs[2:3]))
	fmt.Println("第四个数:", hex.EncodeToString(acs[3:4]))

	return []byte{0xAA, 0xF3, 0x04, 0x11, 0x00, 0x02, 0x61}

}

func UpdatePrice(iccid string, a_price, b_price float64) []byte {
	_ = retuByte(a_price, b_price)

	return []byte{0xaa, 0x00}

}
