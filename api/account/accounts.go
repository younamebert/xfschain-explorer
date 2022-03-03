package api

import (
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

	var (
		page     int
		pageSize int
		result   *apis.Pages = new(apis.Pages)
	)

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = conf.Page
	}

	pageSize, err = strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		pageSize = conf.PageSize
	}
	counts := ac.Handle.HandleChainAddress.Count(nil, nil)
	txs := ac.Handle.HandleChainAddress.GetAccounts(nil, nil, page, pageSize)
	if len(txs) == 0 || txs == nil {
		common.SendResponse(c, http.StatusOK, nil, nil)
		return
	}

	result = &apis.Pages{
		Page:     int64(page),
		PageSize: int64(pageSize),
		Limits:   counts,
		Data:     txs,
	}
	common.SendResponse(c, http.StatusOK, nil, result)
}

func (ac *AccountLinkApi) Detailed(c *gin.Context) {
	var (
		addr   string
		result *apis.Pages = new(apis.Pages)
	)

	addr = c.Query("addr")
	if addr == "" {
		common.SendResponse(c, http.StatusBadRequest, common.NotParamErr, nil)
		return
	}

	addrObj := ac.Handle.HandleChainAddress.Query("address = ?", addr)
	if len(addrObj) == 0 || addrObj == nil {
		common.SendResponse(c, http.StatusOK, nil, nil)
		return
	}
	result = &apis.Pages{
		Data: addrObj[0],
	}
	common.SendResponse(c, http.StatusOK, nil, result)
}

func (ac *AccountLinkApi) DetailedTxs(c *gin.Context) {

	var (
		addr     string
		page     int
		pageSize int
		limits   int64       = 0
		result   *apis.Pages = new(apis.Pages)
	)

	addr = c.Query("addr")
	if addr == "" {
		common.SendResponse(c, http.StatusBadRequest, common.NotParamErr, nil)
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = conf.Page
	}

	pageSize, err = strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		pageSize = conf.PageSize
	}

	txs := ac.Handle.HandleBlockTxs.GetTxs(`tx_from = ?`, addr, page, pageSize)
	if len(txs) == 0 || txs == nil {
		common.SendResponse(c, http.StatusOK, nil, nil)
		return
	}

	limits = ac.Handle.HandleBlockTxs.Count(`tx_from = ?`, addr)
	result = &apis.Pages{
		Page:     int64(page),
		PageSize: int64(pageSize),
		Limits:   limits,
		Data:     txs,
	}
	common.SendResponse(c, http.StatusOK, nil, result)
}
