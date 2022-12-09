package exchange

import (
	"errors"
	"github.com/raozhaofeng/beego/components/exchange/binance"
	"github.com/raozhaofeng/beego/components/exchange/interfaces"
)

const (
	BinanceExchange = "binance"
)

// NewExchange 创建交易所
func NewExchange(exchangeName string) (interfaces.Exchange, error) {
	switch exchangeName {
	case BinanceExchange:
		return binance.NewExchange(), nil
	default:
		return nil, errors.New("NotFoundExchange")
	}
}
