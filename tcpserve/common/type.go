package common

// 功能码
const (
	Header      = byte(0xAA)
	Registers   = byte(0xF0)
	Pant        = byte(0xF1)
	UploadOrder = byte(0xf2)
	SetPrice    = byte(0xf3)
	Switchad    = byte(0xf4)
	SwitchLamp  = byte(0xf5)
)

// 仓库枚举
const (
	EmptyWare = 0 // 空仓 tag:type-0
	HalfWare  = 1 // 半  tag:type-1
	FullWare  = 2 // 满  tag:type-2
)

// 默认空仓
func WareType(ware byte) int {
	switch ware {
	case byte(0x00):
		return EmptyWare
	case byte(0x01):
		return HalfWare
	case byte(0x02):
		return FullWare
	}
	return EmptyWare
}
