package binance

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Resp struct {
	Symbol             string        `json:"symbol"`              //	交易对
	OrderId            int64         `json:"orderId"`             //	订单ID
	OrderListId        int64         `json:"orderListId"`         //	OCO订单的ID，不然就是-1
	ClientOrderId      string        `json:"clientOrderId"`       //	客户自己设置的ID
	TransactTime       time.Duration `json:"transactTime"`        //	交易的时间戳
	Price              string        `json:"price"`               //	订单价格
	OrigQty            string        `json:"origQty"`             //	用户设置的原始订单数量
	ExecutedQty        string        `json:"executedQty"`         //	交易的订单数量
	CumulativeQuoteQty string        `json:"cummulativeQuoteQty"` //	累计交易的金额
	Status             string        `json:"status"`              //	订单状态
	TimeInForce        string        `json:"timeInForce"`         //	订单的时效方式
	Type               string        `json:"type"`                //	订单类型， 比如市价单，现价单等
	Side               string        `json:"side"`                //	订单方向，买还是卖
	Fills              []*RespFills  `json:"fills"`               //	订单中交易的信息
}

// RespFills 订单中的交易信息
type RespFills struct {
	Price           string `json:"price"`           //	交易的价格
	Qty             string `json:"qty"`             //	交易的数量
	Commission      string `json:"commission"`      //	手续费金额
	CommissionAsset string `json:"commissionAsset"` //	手续费的币种
}

// GetFeeQty 获取手续费的数量
func (c *Resp) GetFeeQty() (float64, string) {
	var fee float64
	commissionAsset := "BNB"
	for _, item := range c.Fills {
		if item.CommissionAsset != "BNB" || strings.Contains(c.Symbol, "BNB") {
			commission, _ := strconv.ParseFloat(item.Commission, 64)
			fee += commission
			commissionAsset = item.CommissionAsset
		}
	}
	if fee > 0 {
		fee, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", fee), 64)
	}
	return fee, commissionAsset
}

// GetPrice 获取价格
func (c *Resp) GetPrice() float64 {
	price, _ := strconv.ParseFloat(c.Price, 64)
	return price
}

// GetExecutedQty 获取执行的数量
func (c *Resp) GetExecutedQty() float64 {
	executedQty, _ := strconv.ParseFloat(c.ExecutedQty, 64)
	return executedQty
}

// GetCumulativeQty 获取累计交易的金额
func (c *Resp) GetCumulativeQty() float64 {
	cumulativeQty, _ := strconv.ParseFloat(c.CumulativeQuoteQty, 64)
	return cumulativeQty
}
