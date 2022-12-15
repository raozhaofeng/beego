package interfaces

type Exchange interface {
	// Conf 设置配置文件
	Conf(conf Conf) Exchange

	// TestOrder 测试订单
	TestOrder() Exchange

	// TickerPrice 最新价格
	TickerPrice() (Ticker, error)

	// Buy 买入
	Buy(p map[string]interface{}) (Resp, error)

	// MarketBuy 市价买
	MarketBuy(symbol string, quoteOrderQty float64) (Resp, error)

	// Sell 卖出
	Sell(p map[string]interface{}) (Resp, error)

	// MarketSell 市价卖
	MarketSell(symbol string, quantity float64) (Resp, error)

	// SetProxy 设置代理
	SetProxy(url string) Exchange

	// SetBaseURL 设置基础地址
	SetBaseURL(baseURL string) Exchange
}
