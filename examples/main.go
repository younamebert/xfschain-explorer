package main

import (
	"fmt"
	"mi/pay/wechatpay"
)

func main() {
	//调用付款支付
	wechatpayclient := wechatpay.DefaultWeChatPay()
	resp, err := wechatpayclient.Micropay(1, "付款码", "订单内容或标题", "发送订单服务ip")
	fmt.Printf("resq:%v err:%v", resp, err)
}
