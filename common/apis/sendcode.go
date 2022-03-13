package apis

//SwitchIedOpen-打开灯带：AA F5 01 FF 02 9F
//SwitchIedClose-关闭灯带：AA F5 01 00 01 A0
var (
	switchIedOpen  = []byte{0xAA, 0xF5, 0x01, 0xFF, 0x02, 0x9F}
	switchIedClose = []byte{0xAA, 0xF5, 0x01, 0x00, 0x01, 0xA0}
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
