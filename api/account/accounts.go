package api

import (
	"errors"
	"net/http"
	"strconv"
	"xfschainbrowser/common"
	"xfschainbrowser/common/apis"
	"xfschainbrowser/conf"

	"github.com/gin-gonic/gin"
)

type AccountLinkApi struct {
	Handle *apis.LinkApi
}

func (ac *AccountLinkApi) GetAccounts(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = conf.Page
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		pageSize = conf.PageSize
	}
	result := new(apis.Pages)
	if page > 0 && pageSize > 0 {
		counts := ac.Handle.HandleChainAddress.Count()
		txs := ac.Handle.HandleChainAddress.GetAccounts(nil, nil, page, pageSize)

		result = &apis.Pages{
			Page:     int64(page),
			PageSize: int64(pageSize),
			Limits:   counts,
			Data:     txs,
		}
		common.SendResponse(c, http.StatusOK, nil, result)
	} else {
		common.SendResponse(c, http.StatusOK, nil, result)
	}
}

func (ac *AccountLinkApi) Detailed(c *gin.Context) {
	addr := c.Query("addr")
	err := errors.New("wallet address not nil")
	if addr == "" {
		common.SendResponse(c, http.StatusBadRequest, err, nil)
		return
	}
	addrObj := ac.Handle.HandleChainAddress.Query("address = ?", addr)
	if addrObj == nil {
		common.SendResponse(c, http.StatusBadRequest, err, nil)
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	pageSize := 20

	result := new(DetailedResp)

	result.Account = addrObj[0]
	if page > 0 && pageSize > 0 {
		counts := addrObj[0].TxCount
		txs := ac.Handle.HandleBlockTxs.GetTxs("from = ?", addr, page, pageSize)

		result.TxsOther = &apis.Pages{
			Page:     int64(page),
			PageSize: int64(pageSize),
			Limits:   int64(counts),
			Data:     txs,
		}
		common.SendResponse(c, http.StatusOK, nil, result)
	} else {
		common.SendResponse(c, http.StatusOK, nil, result)
	}
	if err != nil {
		common.SendResponse(c, http.StatusBadRequest, err, nil)
		return
	}
}
