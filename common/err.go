package common

import "errors"

const (
	SystemErr      = "-20001"
	CustomErr      = "-30001"
	MySqlUniqueErr = "1062:"
)

var (
	NotDataErr  = errors.New("nonexistent data")
	NotParamErr = errors.New("parameter cannot be empty")
)

// func ErrCode(err error) error {
// 	ms := strings.Fields(err.Error())
// 	fmt.Printf("ms:%v\n", ms[1])
// 	if ms[1] == MySqlUniqueErr {
// 		return nil
// 	}
// 	return err
// }
