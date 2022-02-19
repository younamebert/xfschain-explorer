package api

import (
	"net/http"
	"xfschainbrowser/common"

	"github.com/gin-gonic/gin"
)

type IndexLinkApi struct{}

func (i *IndexLinkApi) Status(c *gin.Context) {
	common.SendResponse(c, http.StatusOK, nil, "nihao")
}
