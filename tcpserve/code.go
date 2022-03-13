package tcpserve

import "mi/tcpserve/common"

// Header      = byte(0xAA)
// Registers   = byte(0xF0)
// AA F0 01 00 01 9B

// 修改成功：AA F3 01 00 01 9E
// 修改失败：AA F3 01 01 01 9F
var (
	AddEquipmentRegisters    = []byte{common.Header, common.Registers, 0x01, 0x00, 0x01, 0x9b}
	AddEquipmentRegistersErr = []byte{common.Header, common.Registers, 0x01, 0x01, 0x01, 0x9c}

	//修改单价
	// 修改成功：AA F3 01 00 01 9E
	// 修改失败：AA F3 01 01 01 9F
	ModifySuccess = []byte{common.Header, common.SetPrice, 0x01, 0x00, 0x01, 0x9E}
	ModifyError   = []byte{common.Header, common.SetPrice, 0x01, 0x1, 0x01, 0x9F}

	//打开广告屏：AA F4 01 FF 02 9E
	// 关闭广告屏：AA F4 01 00 01 9F

	SwitchAdvertisingSucc  = []byte{common.Header, common.Switchad, 0x01, 0xFF, 0x02, 0x9E}
	SwitchAdvertisingError = []byte{common.Header, common.Switchad, 0x01, 0x00, 0x01, 0x9F}

	// 打开灯带：AA F5 01 FF 02 9F
	// 关闭灯带：AA F5 01 00 01 A0
	SwitchBeltSucc  = []byte{common.Header, common.SwitchLamp, 0x01, 0xFF, 0x02, 0x9F}
	SwitchBeltError = []byte{common.Header, common.SwitchLamp, 0x01, 0x00, 0x01, 0xA0}

	// 成功：AA F1 01 FF 02 9B
	// 失败：AA F1 01 00 02 9A
	HeartbeatSucc  = []byte{common.Header, common.Pant, 0x01, 0xFF, 0x02, 0x9B}
	HeartbeatError = []byte{common.Header, common.Pant, 0x01, 0x00, 0x02, 0x9A}

	//收款成功：AA F2 01 00 01 9D
	//余额不足，收款失败：AA F2 01 01 01 9E
	//付款超时：AA F2 01 02 01 9F
	UploadOrderSucc     = []byte{common.Header, common.UploadOrder, 0x01, 0x00, 0x01, 0x9D}
	UploadOrderError    = []byte{common.Header, common.UploadOrder, 0x01, 0x01, 0x01, 0x9E}
	UploadOrderOvertime = []byte{common.Header, common.UploadOrder, 0x01, 0x02, 0x01, 0x9F}
)
