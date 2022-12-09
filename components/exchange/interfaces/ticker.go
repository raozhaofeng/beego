package interfaces

// Ticker 交易对最新价格
type Ticker interface {
	// SymbolPriceMap 交易对最新价格列表
	SymbolPriceMap() map[string]float64

	// SymbolPrice 交易对最新价格
	SymbolPrice(symbol string) float64

	// TickerListToTickerMap 最新交易对价格数组转最新交易对Map
	TickerListToTickerMap() map[string]float64
}
