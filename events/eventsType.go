package events

type SendNoticeEvent struct {
	Iccid string
	Data  []byte
}
