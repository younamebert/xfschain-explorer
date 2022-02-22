package common

import (
	"fmt"
	"strings"
)

const (
	SystemErr      = "-20001"
	CustomErr      = "-30001"
	MySqlUniqueErr = "1062:"
)

func ErrCode(err error) error {
	ms := strings.Fields(err.Error())
	fmt.Printf("ms:%v\n", ms[1])
	if ms[1] == MySqlUniqueErr {
		return nil
	}
	return err
}
