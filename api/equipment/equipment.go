package api

import (
	"fmt"
	"mi/common"
	"mi/common/apis"
	"mi/conf"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EquipmentLinkApi struct {
	Handle *apis.LinkApi
}

func (ac *EquipmentLinkApi) EquipmentList(c *gin.Context) {
	var (
		page     int         = conf.Page
		pageSize int         = conf.PageSize
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
	count := ac.Handle.HandleMiEquipment.Count(nil, nil)
	fmt.Println(count)
	fmt.Println(page, pageSize)
	list := ac.Handle.HandleMiEquipment.Querys(nil, nil, page, pageSize)
	fmt.Println(list)
	if len(list) < 1 {
		common.SendResponse(c, http.StatusOK, nil, nil)
		return
	}

	result = &apis.Pages{
		Page:     int64(page),
		PageSize: int64(pageSize),
		Limits:   count,
		Data:     list,
	}
	common.SendResponse(c, http.StatusOK, nil, result)
}
func (ac *EquipmentLinkApi) EquipmentSwitch(c *gin.Context) {
	//   c.Qquer
	// ac.Handle.HandleMiEquipment.Insert()
}

//关闭/打开广告屏 0关 1开
func (ac *EquipmentLinkApi) SwitchAdvertising(c *gin.Context) {
	//设备id
	iccid := c.Query("iccid")

	status := c.Query("status")

	if iccid == "" || status == "" {
		return
	}

	types, _ := strconv.Atoi(status)

	err := ac.Handle.HandleMiEquipment.SetSwitchad(iccid, types)

	if err != nil {

		common.SendResponse(c, http.StatusOK, nil, "失败")
	}

	//事件

	common.SendResponse(c, http.StatusOK, nil, "成功")
}

//关闭/打开灯带0关 1开
func (ac *EquipmentLinkApi) SwitchLed(c *gin.Context) {
	//iccid
	iccid := c.Query("iccid")

	status := c.Query("status")

	types, _ := strconv.Atoi(status)

	err := ac.Handle.HandleMiEquipment.SetSwitchadLed(iccid, types)

	if err != nil {
		common.SendResponse(c, http.StatusOK, nil, "失败")
	}
	common.SendResponse(c, http.StatusOK, nil, "成功")
}
