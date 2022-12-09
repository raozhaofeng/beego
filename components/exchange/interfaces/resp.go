package interfaces

type Resp interface {
	// GetPrice 订单单价
	GetPrice() float64

	// GetExecutedQty 获取执行的数量
	GetExecutedQty() float64

	// GetCumulativeQty 获取累计交易的金额
	GetCumulativeQty() float64

	// GetFeeQty 获取手续费的数量
	GetFeeQty() (float64, string)
}
