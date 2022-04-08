package wechatpay

import (
	"context"
	"fmt"
	"mi/conf"
	"mi/global"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/util"
	"github.com/go-pay/gopay/wechat"
)

type WeChatPayConfig struct {
	APPID       string `json:"appid"`
	MchId       string `json:"mchid"`
	APIv2Key    string `json:"apiv2key"`
	IsProd      bool   `json:"is_prod"`
	DebugSwitch int    `json:"debug_switch"`
}

func defaultWeChatPayConfig() *WeChatPayConfig {
	return &WeChatPayConfig{
		APPID:       conf.APPID,
		MchId:       conf.MchId,
		APIv2Key:    conf.APIv2Key,
		IsProd:      conf.IsProd,
		DebugSwitch: conf.DebugOn,
	}
}

type WeChatPay struct {
	config       *WeChatPayConfig
	WeChatClient *wechat.Client
}

func NewWeChatPay(conf *WeChatPayConfig) *WeChatPay {
	result := new(WeChatPay)
	result.config = conf
	result.WeChatClient = wechat.NewClient(result.config.APPID, result.config.MchId, result.config.APIv2Key, result.config.IsProd)
	result.WeChatClient.DebugSwitch = gopay.DebugSwitch(result.config.DebugSwitch)
	return result
}

func DefaultWeChatPay() *WeChatPay {
	result := new(WeChatPay)
	result.config = defaultWeChatPayConfig()
	result.WeChatClient = wechat.NewClient(result.config.APPID, result.config.MchId, result.config.APIv2Key, result.config.IsProd)
	result.WeChatClient.DebugSwitch = gopay.DebugSwitch(result.config.DebugSwitch)
	return result
}

/*
title:Micropay
description:通过付款码来进行微信支付
auth:schoolbert@foxmail.com 22/4/7
param:
[
	{name:body,type:string,des:"支付订单标题"},
	{name:auth_code,type:string,des:"付款码"},
	{name:total_fee,type:string,des:"付款金额"}, // 1 == 0.01元RMB
	{name:spbill_create_ip,type:string,des:"订单创建服务器IP"}
]
return:
[
	{type:wechat.MicropayResponse},
	{type:error}
]
*/
func (wechatpay *WeChatPay) Micropay(total_fee int, auth_code, body, spbill_create_ip string) (*wechat.MicropayResponse, error) {
	bm := make(gopay.BodyMap)

	orderNumber := util.RandomString(32)
	bm.Set("nonce_str", util.RandomString(32)).
		Set("body", body).
		Set("out_trade_no", orderNumber).
		Set("total_fee", total_fee).
		Set("spbill_create_ip", spbill_create_ip).
		Set("auth_code", auth_code).
		Set("sign_type", wechat.SignType_MD5)

	wxRsp, err := wechatpay.WeChatClient.Micropay(context.Background(), bm)
	if err != nil {
		return wxRsp, err
	}

	ok, err := wechat.VerifySign(wechatpay.config.APIv2Key, wechat.SignType_MD5, wxRsp)
	if err != nil {
		return nil, err
	}
	global.GVA_LOG.Info(fmt.Sprintf("SignOk:%v", ok))
	return wxRsp, nil
}

/*
title:OrderQuery
description:查询订单付款信息
auth:schoolbert@foxmail.com 22/4/7
param:
[
	{name:out_trade_no,type:string,des:"订单号"},
]
return:
[
	{type:wechat.QueryOrderResponse},
	{type:gopay.BodyMap},
	{type:error}
]
*/
func (wechatpay *WeChatPay) OrderQuery(out_trade_no string) (*wechat.QueryOrderResponse, gopay.BodyMap, error) {

	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", out_trade_no).
		Set("nonce_str", util.RandomString(32)).
		Set("sign_type", wechat.SignType_MD5)

	return wechatpay.WeChatClient.QueryOrder(context.Background(), bm)
}
