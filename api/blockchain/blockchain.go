package api

import (
	"mi/common/apis"

	"github.com/gin-gonic/gin"
)

type BlocksLinkApi struct {
	Handle *apis.LinkApi
}

func (bc *BlocksLinkApi) GetBlocks(c *gin.Context) {

}

func (bc *BlocksLinkApi) Detailed(c *gin.Context) {

}

func (bc *BlocksLinkApi) DetailedTx(c *gin.Context) {

}
