package tcpserve

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"mi/global"
	"mi/model"
	"mi/tcpserve/common"
	"mi/tools"
	"net"

	"github.com/shopspring/decimal"
)

type Handle struct {
	model *model.RecordHandle
	iccid string
}

func NewHandle() *Handle {
	return &Handle{
		model: model.NewRecordHandle(),
	}
}

// func (h *Handle)
func (h *Handle) Process(conn net.Conn) error {

	//打开注释的地方(devour)
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		buf := make([]byte, 128)
		n, err := reader.Read(buf) // 读取数据
		if err != nil {
			global.GVA_LOG.Warn(err.Error())
			return err
		}
		result := buf[:n]
		write, err := h.chck(result)
		if err != nil {
			global.GVA_LOG.Warn(err.Error())
			return err
		}
		if _, err := conn.Write(write); err != nil {
			global.GVA_LOG.Warn(err.Error())
			return fmt.Errorf("write client msg, err: %v", err)
		}
		// fmt.Printf("%x", result)
		// fmt.Println(23456)
	}
	return nil
}

// func (h *Handle)
func (h *Handle) chck(data []byte) ([]byte, error) {
	//验证数据是否AA开头
	if data[0] != common.Header {
		return nil, errors.New("header data error")
	}

	//验证crc16
	crc16 := common.CRC(data[:len(data)-2])
	if crc16 != fmt.Sprintf("%X", data[len(data)-2:]) {
		return nil, errors.New("header crc error")
	}

	//验证数据长度
	dataLen := data[2:3]
	pending := data[3 : len(data)-2]
	if int(dataLen[0]) != len(pending) {
		return nil, errors.New("header crc error")
	}

	// if int(data[2])
	// 优化
	result, err := h.work(data[1], pending)
	if err != nil {
		return nil, err
	}
	return result, nil

	// if _, err := h.Conn.Write(result); err != nil {
	// 	global.GVA_LOG.Warn(fmt.Sprintf("iccid:%v session reply error:%v", h.Iccid, err))
	// 	h.Conn.Close()
	// 	return
	// }
}

// arrary("nihao"=>["nihao":1])
func (h *Handle) work(funccode byte, data []byte) ([]byte, error) {
	switch funccode {
	case common.Registers:
		return h.registers(data)
	case common.Pant:
		return h.pant(data)
	case common.UploadOrder:
		return h.uploadOrder(data)
	case common.SetPrice:
		return h.setPrice(data)
	case common.Switchad:
		return h.switchad(data)
	case common.SwitchLamp:
		return h.switchLamp(data)
	}

	return nil, nil
}

func (h *Handle) registers(data []byte) ([]byte, error) {
	iccdRune := make([]rune, 0)
	for _, v := range data {
		iccdRune = append(iccdRune, rune(v))
	}
	//Ascii码
	iccid := string(iccdRune)

	write := &model.MiEquipment{
		Iccid: iccid,
	}
	if err := h.model.HandleMiEquipment.Insert(write); err != nil {
		fmt.Println(err)
		return AddEquipmentRegistersErr, err
	}
	h.iccid = iccid
	// 写入数据库
	return AddEquipmentRegisters, nil
}

func (h *Handle) pant(data []byte) ([]byte, error) {
	var (
		wareA int = 0 // A仓库
		wareB int = 0 // B仓库
	)
	fmt.Println(h.iccid)
	wareA = common.WareType(data[0])
	wareB = common.WareType(data[1])
	// 写入数据库(a,b仓库)
	a := wareA
	b := wareB

	//查询数据库存不存在

	EmptyA := h.model.HandleMiWarehouse.QueryOne("iccid =?", h.iccid, "status=?", 1)

	//
	//查询设备的a,b仓价格
	list := h.model.HandleMiEquipment.QueryOne("iccid =?", h.iccid)

	fmt.Println("数据:", list.AWarehousePrice)
	table1 := &model.MiWarehouse{
		Iccid:          h.iccid,
		WarehouseType:  a,
		WarehousePrice: list.AWarehousePrice,
		Status:         1,
	}
	if EmptyA.Iccid == "" {
		h.model.HandleMiWarehouse.Insert(table1)

	} else {
		//修改

		EmptyA.WarehousePrice = list.AWarehousePrice
		EmptyA.Status = b
		h.model.HandleMiWarehouse.SaveWare("iccid =?", h.iccid, "status=?", 1, EmptyA)
	}

	EmptyB := h.model.HandleMiWarehouse.QueryOne("iccid =?", h.iccid, "status=?", 2)

	table2 := &model.MiWarehouse{
		Iccid:          h.iccid,
		WarehouseType:  b,
		WarehousePrice: list.BWarehousePrice,
		Status:         2,
	}
	if EmptyB.Iccid == "" {
		h.model.HandleMiWarehouse.Insert(table2)
	} else {
		//修改

		EmptyB.WarehousePrice = list.BWarehousePrice
		EmptyB.Status = b
		h.model.HandleMiWarehouse.SaveWare("iccid =?", h.iccid, "status=?", 2, EmptyB)
	}

	mangeList := h.model.HandleWlMange.Query("iccid =?", h.iccid)

	//a仓位预警
	earlyA := h.model.HandleWlEarly.QueryOne("iccid =?", h.iccid, "warehouse=?", 1)

	if earlyA.Iccid == "" {

		earlyA_inert := &model.WlEarly{
			MangeId:   int(mangeList.MangeId),
			Number:    h.iccid,
			Type:      1,
			CStatus:   a,
			Warehouse: 1,
			Iccid:     h.iccid,
		}

		h.model.HandleWlEarly.Insert(earlyA_inert)
	} else {
		earlyA.CStatus = a
		h.model.HandleWlEarly.SaveEarly("iccid =?", h.iccid, "warehouse=?", 1, earlyA)
	}

	//b仓库预警
	earlyB := h.model.HandleWlEarly.QueryOne("iccid =?", h.iccid, "warehouse=?", 2)

	if earlyB.Iccid == "" {

		earlyB_inert := &model.WlEarly{
			MangeId:   int(mangeList.MangeId),
			Number:    h.iccid,
			Type:      1,
			CStatus:   b,
			Warehouse: 2,
			Iccid:     h.iccid,
		}

		h.model.HandleWlEarly.Insert(earlyB_inert)
	} else {
		earlyB.CStatus = b
		h.model.HandleWlEarly.SaveEarly("iccid =?", h.iccid, "warehouse=?", 2, earlyB)

	}

	return HeartbeatSucc, nil
}

//上传订单
func (h *Handle) uploadOrder(data []byte) ([]byte, error) {
	var (
		amounts decimal.Decimal
		payType int
		payCode string
	)

	moneyUint64 := common.Hex2int(&[]byte{data[0], data[1]}) // 1045
	amounts = common.Uint64toDecimal(int64(moneyUint64), 100)

	fmt.Println(amounts, payType, payCode)

	cardNumber := h.model.HandleWlCardNumber.Query("number =?", "8939131")

	//余额不足
	money := decimal.NewFromFloat(cardNumber.Money).Round(2).BigFloat()

	if money.Cmp(amounts.BigFloat()) < 0 {
		return UploadOrderError, nil
	}

	write := &model.MiOrder{
		Iccid:         h.iccid,
		PayType:       len(data),
		PayCode:       strconv.FormatInt(time.Now().Unix()+int64(rand.Intn(100000)), 10),
		OrderNumber:   tools.GetOrder(),
		PaymentAmount: amounts,
		Number:        "8939131",
	}

	//查询设备id
	list := h.model.HandleWlMange.Query("number =?", h.iccid)

	amouts, _ := amounts.Float64()

	//设备出售米纪录表
	saleWari := &model.WlSale{
		MangeId: list.MangeId,
		Number:  2,
		Money:   amouts,
	}

	//购米记录表
	writes := &model.WlPurchaseRecord{
		Price:    amouts,
		OrderNo:  tools.GetOrder(),
		PayTime:  time.Now(),
		MemberId: cardNumber.MemberId,
		Iccid:    h.iccid,
		MangeId:  list.MangeId,
	}

	if money.Cmp(amounts.BigFloat()) > 0 {
		if err := h.model.HandleMiOrder.Insert(write); err != nil {
			global.GVA_LOG.Warn(err.Error())

			//返回错误码
			return AddEquipmentRegistersErr, err
		}

		//设备出售米纪录表
		if err := h.model.HandleWlSale.Insert(saleWari); err != nil {
			global.GVA_LOG.Warn(err.Error())
		}

		//加入购米记录表
		h.model.HandleWlPurchaseRecord.Insert(writes)

		if err := h.model.HandleWlSale.Insert(saleWari); err != nil {
			global.GVA_LOG.Warn(err.Error())
		}

		//扣余额
		money.Sub(money, amounts.BigFloat())

		moneyFloat, _ := money.Float64()

		cardNumber.Money = moneyFloat

		h.model.HandleWlCardNumber.Update("card_number_id=?", cardNumber.CardNumberId, cardNumber)

		return UploadOrderSucc, nil
	}

	if money.Cmp(amounts.BigFloat()) == 0 {
		if err := h.model.HandleMiOrder.Insert(write); err != nil {
			// 	//返回错误码
			global.GVA_LOG.Warn(err.Error())
			return AddEquipmentRegistersErr, err
		}

		//加入购米记录表
		if err := h.model.HandleWlPurchaseRecord.Insert(writes); err != nil {
			global.GVA_LOG.Warn(err.Error())
		}

		//设备出售米纪录表
		h.model.HandleWlSale.Insert(saleWari)
		//扣余额
		// money.Sub(money, amounts.BigFloat())
		// h.model.HandleWlCardNumber.Update("number =?", "8939131", "money", moneys)

		moneyFloat, _ := money.Float64()
		cardNumber.Money = moneyFloat
		h.model.HandleWlCardNumber.Update("card_number_id=?", cardNumber.CardNumberId, cardNumber)

		return UploadOrderSucc, nil
	}

	// switch len(data[1:]) {
	// case 10:

	// 	// if err := h.model.HandleMiOrder.Insert(write); err != nil {
	// 	// 	//返回错误码
	// 	// 	fmt.Println("错误:", err)
	// 	// 	return AddEquipmentRegistersErr, err
	// 	// }
	// case 18:

	// }

	return AddEquipmentRegistersErr, nil
}

func (h *Handle) setPrice(data []byte) ([]byte, error) {
	// var (
	// 	wareAPrice decimal.Decimal
	// 	wareBPrice decimal.Decimal
	// )

	// wareAPriceUint64 := common.Hex2int(&[]byte{data[0], data[1]}) //350

	// wareAPrice = common.Uint64toDecimal(int64(wareAPriceUint64), 100) //3.50

	// wareBPriceUint64 := common.Hex2int(&[]byte{data[2], data[3]}) //
	// wareBPrice = common.Uint64toDecimal(int64(wareBPriceUint64), 100)

	// // fmt.Printf("a:%v b:%v\n", wareAPrice.String(), wareBPrice.String())
	// whereA := make([]interface{}, 0)
	// whereA = append(whereA, h.Iccid, 1)
	// warehouseA := h.model.HandleMiWarehouse.Query("iccid = ? and warehouse_type = ?", whereA)
	// if len(warehouseA) > 0 {
	// 	warehouseAWrite := warehouseA[0]
	// 	rPrice, _ := wareAPrice.Round(2).Float64()
	// 	warehouseAWrite.WarehousePrice = rPrice
	// 	if err := h.model.HandleMiWarehouse.Update(warehouseAWrite); err != nil {
	// 		//错误的码
	// 		// common.ModifyError
	// 		return ModifyError, nil
	// 	}
	// } else {
	// 	warehouseAWrite := &model.MiWarehouse{
	// 		Iccid:          h.Iccid,
	// 		WarehouseType:  1,
	// 		WarehousePrice: float64(0),
	// 	}
	// 	if err := h.model.HandleMiWarehouse.Insert(warehouseAWrite); err != nil {
	// 		//错误的码
	// 		return ModifyError, nil
	// 	}
	// }

	// whereB := make([]interface{}, 0)
	// whereB = append(whereB, h.Iccid, 0)
	// warehouseB := h.model.HandleMiWarehouse.Query("iccid = ? and warehouse_type = ?", whereB)
	// if len(warehouseB) > 0 {
	// 	warehouseBWrite := warehouseB[0]
	// 	rPrice, _ := wareBPrice.Round(2).Float64()
	// 	warehouseBWrite.WarehousePrice = rPrice
	// 	if err := h.model.HandleMiWarehouse.Update(warehouseBWrite); err != nil {
	// 		//错误的码
	// 		return ModifyError, nil
	// 	}
	// } else {
	// warehouseBWrite := &model.MiWarehouse{
	// 		Iccid:          h.Iccid,
	// 		WarehouseType:  1,
	// 		WarehousePrice: float64(0),
	// 	}
	// 	if err := h.model.HandleMiWarehouse.Insert(warehouseBWrite); err != nil {
	// 		//错误的码
	// 		return ModifyError, nil
	// 	}
	// }
	return nil, nil
}

func (h *Handle) switchad(data []byte) ([]byte, error) {
	var switchadType int = 0
	switch data[0] {
	case byte(0x00):
		switchadType = 1
	case byte(0xFF):
		switchadType = 0
	}
	fmt.Println(switchadType)
	return SwitchAdvertisingSucc, nil
}

func (h *Handle) switchLamp(data []byte) ([]byte, error) {
	var switchLampType int = 0
	switch data[0] {
	case byte(0x00):
		switchLampType = 1
	case byte(0xFF):
		switchLampType = 0
	}
	fmt.Println(switchLampType)
	return SwitchBeltSucc, nil
}
