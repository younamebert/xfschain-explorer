package api

import (
	"mi/common/apis"

	"github.com/gin-gonic/gin"
)

type TxsLinkApi struct {
	Handle *apis.LinkApi
}

func (tx *TxsLinkApi) GetTxs(c *gin.Context) {

}

func (tx *TxsLinkApi) Detailed(c *gin.Context) {
}
