package binance

import "time"

// SpecificationInfo 交易所规范信息
var SpecificationInfo = &Info{}

// SpecificationSymbolInfo 交易所交易对规范
var SpecificationSymbolInfo = map[string]*InfoSymbol{}

type Info struct {
	Timezone   string        `json:"timezone"`   //	时区
	ServerTime time.Duration `json:"serverTime"` //	服务器时间
	//	交易平台API使用的一些次数要求
	//	rateLimits有三种类型，别分是
	//		一分钟内请求权重之和的上限
	//		每秒钟交易次数的上限
	//		每天交易次数的上限
	RateLimits []*InfoRateLimits `json:"rateLimits"`
	Symbols    []*InfoSymbol     `json:"symbols"`
}

// InfoRateLimits api 要求
type InfoRateLimits struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

// InfoSymbol 交易对
type InfoSymbol struct {
	Symbol                     string                            `json:"symbol"`                     //	交易对
	Status                     string                            `json:"status"`                     //	状态，trading应该表示可交易
	BaseAsset                  string                            `json:"baseAsset"`                  //
	BaseAssetPrecision         int                               `json:"baseAssetPrecision"`         //	资产精确的小数点位数，这里是8位
	QuoteAsset                 string                            `json:"quoteAsset"`                 //
	QuotePrecision             int                               `json:"quotePrecision"`             //
	QuoteAssetPrecision        int                               `json:"quoteAssetPrecision"`        //	同理
	BaseCommissionPrecision    int                               `json:"baseCommissionPrecision"`    //
	QuoteCommissionPrecision   int                               `json:"quoteCommissionPrecision"`   //
	OrderTypes                 []string                          `json:"orderTypes"`                 //	交易的类别有哪些
	IcebergAllowed             bool                              `json:"icebergAllowed"`             //	是否允许——冰山订单
	OcoAllowed                 bool                              `json:"ocoAllowed"`                 //	是否允许——OCO订单
	QuoteOrderQtyMarketAllowed bool                              `json:"quoteOrderQtyMarketAllowed"` //
	IsSpotTradingAllowed       bool                              `json:"isSpotTradingAllowed"`       //	是否允许——现货交易订单
	IsMarginTradingAllowed     bool                              `json:"isMarginTradingAllowed"`     //	是否允许——保证金交易订单
	Filters                    []map[string]interface{}          `json:"filters"`                    //	筛选器，订单不满足这些筛选器的要求是无法进行交易的
	FiltersMap                 map[string]map[string]interface{} `json:"filtersMap"`                 //	筛选器，处理之后的
}

// InitSpecification 初始化交易所规范
func InitSpecification() {
	specification, err := NewExchange().SpecificationInfo()
	if err != nil {
		panic(err)
	}

	SpecificationInfo = specification
	var symbolInfoMap = map[string]*InfoSymbol{}
	for _, symbolInfo := range specification.Symbols {
		// 筛选器，转Map
		for _, symbolFilters := range symbolInfo.Filters {
			symbolInfo.FiltersMap[symbolFilters["filterType"].(string)] = symbolFilters
		}
		symbolInfoMap[symbolInfo.Symbol] = symbolInfo
	}
	SpecificationSymbolInfo = symbolInfoMap
}
