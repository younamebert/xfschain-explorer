package api

import (
	"net/http"
	"strconv"
	"xfschainbrowser/common"
	"xfschainbrowser/common/apis"
	"xfschainbrowser/conf"

	"github.com/gin-gonic/gin"
)

type TxsLinkApi struct {
	Handle *apis.LinkApi
}

func (tx *TxsLinkApi) GetTxs(c *gin.Context) {

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

	counts := tx.Handle.HandleBlockTxs.Count(nil, nil)
	txs := tx.Handle.HandleBlockTxs.GetTxs(nil, nil, page, pageSize)
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

func (tx *TxsLinkApi) Detailed(c *gin.Context) {
	var (
		txhash string
		result *apis.Pages = new(apis.Pages)
	)
	txhash = c.Query("hash")
	txs := tx.Handle.HandleBlockTxs.Query("hash = ?", txhash)
	if len(txs) == 0 || txs == nil {
		common.SendResponse(c, http.StatusOK, nil, nil)
		return
	}
	result = &apis.Pages{
		Data: txs[0],
	}
	common.SendResponse(c, http.StatusOK, nil, result)
}
