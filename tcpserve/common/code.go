package common

import "net"

const (
	CardPay   = 1  //卡片支付
	WeChatPay = 2  // 微信支付
	Alipay    = 3  //支付宝支付
	Card      = 10 //刷卡支付
	Other     = 18 //微信和支付支付
)

func GetLocalIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		return ipAddr.IP.String(), nil
	}
	return
}
