package api

import (
	"errors"
	"fmt"
	"mi/common"
	"mi/common/apis"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay/wechat"
)

type PayLinkApi struct {
	Handle *apis.LinkApi
}

func (pay *PayLinkApi) WeChatNotify(c *gin.Context) {

	wxnotify, err := wechat.ParseNotify(c.Request)
	if err != nil {
		common.SendResponse(c, http.StatusBadRequest, err, nil)
		return
	}
	if wxnotify.ResultCode != "SUCCESS" && wxnotify.ReturnCode != "SUCCESS" {
		common.SendResponse(c, http.StatusBadRequest, errors.New(wxnotify.ReturnMsg), nil)
		return
	}
	ordernumber := wxnotify.OutTradeNo
	fmt.Printf("orderNumber:%v\n", ordernumber)
	// wechat.VerifySign(conf.APIv2Key, wechat.SignType_MD5, wxnotify)
}
