package apis

//SwitchIedOpen-打开灯带：AA F5 01 FF 02 9F
//SwitchIedClose-关闭灯带：AA F5 01 00 01 A0
var (
	switchIedOpen  = []byte{0xAA, 0xF5, 0x01, 0xFF, 0x02, 0x9F}
	switchIedClose = []byte{0xAA, 0xF5, 0x01, 0x00, 0x01, 0xA0}
)

// 服务端发送：
// 打开广告屏：AA F4 01 FF 02 9E
// 关闭广告屏：AA F4 01 00 01 9F

var (
	screenOpen  = []byte{0xAA, 0xF4, 0x01, 0xFF, 0x02, 0x9E}
	screenClose = []byte{0xAA, 0xF4, 0x01, 0x00, 0x01, 0x9F}
)

//开
// aaf6034951550299

//关
//aaf6034950510294
var (
	SwitchOpen  = []byte{0xAA, 0xF6, 0x03, 0x49, 0x51, 0x55, 0x02, 0x92}
	SwitchClose = []byte{0xAA, 0xF6, 0x03, 0x49, 0x50, 0x51, 0x02, 0x94}
)

const (
	CloseStatus int = 0
	OpenStatus  int = 1
)

func SwitchIed(status int) []byte {
	if status == 1 {
		return switchIedOpen
	} else {
		return switchIedClose
	}
}

func Screen(status int) []byte {
	if status == 1 {
		return screenOpen
	} else {
		return screenClose
	}
}

//开关机
func SwitchMac(status int) []byte {
	if status == 1 {
		return SwitchOpen
	} else {
		return SwitchClose
	}
}
