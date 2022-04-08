package tcpserve

import (
	"fmt"
	"mi/events"
	"mi/model"
	"mi/pay/wechatpay"
	"testing"
)

func setupHandle() *Handle {
	handle := &Handle{
		model:          model.NewRecordHandle(),
		eventBus:       events.EventBusExample,
		wechatpayServe: wechatpay.DefaultWeChatPay(),
		iccid:          "131375319582223286", //前提数据库要有这条数据
	}
	return handle
}

func TestWechatPay(t *testing.T) {
	handle := setupHandle()
	auto_code := "132627763299749640"
	amount := 1
	result, err := handle.WechatPay(amount, auto_code)
	fmt.Printf("%v err:%v", result, err)
}
