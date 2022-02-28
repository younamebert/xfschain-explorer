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
		counts := tx.Handle.HandleBlockTxs.Count()
		txs := tx.Handle.HandleBlockTxs.GetTxs(nil, nil, page, pageSize)

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

// func (t *TxsLinkApi)
