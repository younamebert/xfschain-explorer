package api

import (
	"mi/common"
	"mi/common/apis"
	"mi/conf"
	"mi/events"
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
	list := ac.Handle.HandleMiEquipment.Querys(nil, nil, page, pageSize)
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
	status, err := strconv.Atoi(c.Query("status"))

	if iccid == "" || err != nil {
		common.SendResponse(c, http.StatusBadRequest, common.NotParamErr, nil)
		return
	}

	if (status != apis.OpenStatus) && (status != apis.CloseStatus) {
		common.SendResponse(c, http.StatusBadRequest, common.NotParamErr, nil)
		return
	}

	if err := ac.Handle.HandleMiEquipment.SetSwitchad(iccid, status); err != nil {
		common.SendResponse(c, http.StatusBadRequest, err, nil)
	}

	// push events sub	   //创建事件
	ac.Handle.EventsBus.Publish(events.SendNoticeEvent{Iccid: iccid, Data: apis.Screen(status)})

	common.SendResponse(c, http.StatusOK, nil, nil)
}

//关闭/打开灯带0关 1开
func (ac *EquipmentLinkApi) SwitchLed(c *gin.Context) {
	//iccid
	iccid := c.Query("iccid")
	status, err := strconv.Atoi(c.Query("status"))

	if iccid == "" || err != nil {
		common.SendResponse(c, http.StatusBadRequest, common.NotParamErr, nil)
		return
	}

	if (status != apis.OpenStatus) && (status != apis.CloseStatus) {
		common.SendResponse(c, http.StatusBadRequest, common.NotParamErr, nil)
		return
	}

	if err := ac.Handle.HandleMiEquipment.SetSwitchadLed(iccid, status); err != nil {
		common.SendResponse(c, http.StatusBadRequest, err, nil)
	}

	// push events sub	   //创建事件
	ac.Handle.EventsBus.Publish(events.SendNoticeEvent{Iccid: iccid, Data: apis.SwitchIed(status)})
	common.SendResponse(c, http.StatusOK, nil, nil)
}

type Data struct {
	Iccid  string  `json:"iccid"`
	APrice float64 `json:"a_price"`
	BPrice float64 `json:"b_price"`
}

//修改单价
func (ac *EquipmentLinkApi) UpdatePrice(c *gin.Context) {
	json := Data{}

	c.Bind(&json)

	//修改数据库a仓 b仓
	ac.Handle.HandleMiEquipment.Update(json.Iccid, "a_warehouse_price", json.APrice)

	ac.Handle.HandleMiEquipment.Update(json.Iccid, "b_warehouse_price", json.BPrice)

	apis.UpdatePrice(json.Iccid, json.APrice, json.BPrice)

	common.SendResponse(c, http.StatusOK, nil, json)

}
