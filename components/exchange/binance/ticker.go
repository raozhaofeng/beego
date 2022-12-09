package binance

import "strconv"

type ExchangeTicker struct {
	Symbol string `json:"symbol"` //	交易对
	Price  string `json:"price"`  //	最新价格
}

type Ticker struct {
	TickerList []*ExchangeTicker
	TickerMap  map[string]float64
}

func NewTicker() *Ticker {
	return &Ticker{
		make([]*ExchangeTicker, 0),
		make(map[string]float64, 0),
	}
}

// SymbolPriceMap 交易对列表最新价格
func (c *Ticker) SymbolPriceMap() map[string]float64 {
	if len(c.TickerMap) == 0 {
		return c.TickerListToTickerMap()
	}
	return c.TickerMap
}

// SymbolPrice 交易对最新的价格
func (c *Ticker) SymbolPrice(symbol string) float64 {
	if len(c.TickerMap) == 0 {
		c.TickerListToTickerMap()
	}

	if _, ok := c.TickerMap[symbol]; ok {
		return c.TickerMap[symbol]
	}
	return 0
}

// TickerListToTickerMap 最新价格数组转最新价格Map
func (c *Ticker) TickerListToTickerMap() map[string]float64 {
	for _, v := range c.TickerList {
		priceTmp, _ := strconv.ParseFloat(v.Price, 64)
		c.TickerMap[v.Symbol] = priceTmp
	}
	return c.TickerMap
}
