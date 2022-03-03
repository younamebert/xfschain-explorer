package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

//SendResponse 返回数据结构体
func SendResponse(c *gin.Context, code int, err error, data interface{}) {
	var errStr string
	if err != nil {
		errStr = err.Error()
	} else {
		errStr = "success"
	}

	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: errStr,
		Result:  data,
	})
}
