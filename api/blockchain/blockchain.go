package api

import (
	"net/http"
	"strconv"
	"xfschainbrowser/common"
	"xfschainbrowser/common/apis"
	"xfschainbrowser/conf"

	"github.com/gin-gonic/gin"
)

type BlocksLinkApi struct {
	Handle *apis.LinkApi
}

func (bc *BlocksLinkApi) GetBlocks(c *gin.Context) {
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
		counts := bc.Handle.HandleBlockHeader.Count()
		txs := bc.Handle.HandleBlockHeader.GetBlocks(nil, nil, page, pageSize)

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
