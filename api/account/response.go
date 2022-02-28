package api

import (
	"xfschainbrowser/common/apis"
	"xfschainbrowser/model"
)

type DetailedResp struct {
	TxsOther *apis.Pages         `json:"txs_other"`
	Account  *model.ChainAddress `json:"account"`
}
