package tcpserve

import (
	"fmt"
	"mi/events"
	"mi/global"
	"mi/initialize"
	"mi/model"
	"mi/pay/wechatpay"
	"regexp"
	"testing"
)

func setupHandle() *Handle {

	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	initialize.DBList()
	handle := &Handle{
		model:          model.NewRecordHandle(),
		eventBus:       events.EventBusExample,
		wechatpayServe: wechatpay.DefaultWeChatPay(),
		iccid:          "131375319582223286", //前提数据库要有这条数据
	}
	return handle
}

//测试微信支付
func TestWechatPay(t *testing.T) {
	handle := setupHandle()
	auto_code := "132627763299749640"
	amount := 1 // == 0.01
	result, err := handle.WechatPay(amount, auto_code)
	fmt.Printf("%s err:%v", result, err)
}

//验证是否微信付款码
func TestAutoCodeRule(t *testing.T) {
	auto_code := "132627763299749640"
	found, err := regexp.MatchString(`^1[0-5]\d{16}$`, auto_code)
	t.Fatalf("found:%v err:%v", found, err)
}

//测试余额支付
func TestBalancePay(t *testing.T) {
	handle := setupHandle()
	auto_code := "xxx"  //刷卡支付的卡号
	amount := uint64(1) // == 0.01
	result, err := handle.BalancePay(amount, auto_code)
	fmt.Printf("%s err:%v", result, err)
}
