package common

import (
	"math/rand"
	"strconv"
	"time"
)

func GetOrder() string {
	year, month, day := time.Now().Date()

	hour, min, sec := time.Now().Clock()

	number := rand.Intn(100000000)

	order_no := strconv.Itoa((year)) + strconv.Itoa(int(month)) + strconv.Itoa(day) + strconv.Itoa(hour) + strconv.Itoa(min) + strconv.Itoa(sec) + strconv.Itoa(number)

	return order_no
}
