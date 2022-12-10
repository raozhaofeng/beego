package binance

import (
	"fmt"
	"github.com/raozhaofeng/beego/components/exchange/interfaces"
	"github.com/raozhaofeng/beego/utils"
	"time"
)

// BroadcastMessageTypeTickerPrice 广播最新价格类型
const BroadcastMessageTypeTickerPrice = "TickerPrice"

// RealTimeTickerPrice 实时最新价格
var RealTimeTickerPrice interfaces.Ticker

// BroadcastService 广播服务
var BroadcastService *Broadcast

// Broadcast 广播
type Broadcast struct {
	Service  *utils.Broadcast //	广播服务
	exchange *Exchange        //	交易所对象
}

// InitBroadcastService 初始化广播服务
func InitBroadcastService() {
	// 初始化广播
	broadcastService := utils.NewBroadcast()
	BroadcastService = &Broadcast{
		Service:  broadcastService,
		exchange: NewExchange(),
	}
}

// SyncTickerPrice 同步最新价格
func (c *Broadcast) SyncTickerPrice(second time.Duration) {
	go func() {
		for {
			// 请求数据
			ticker, err := c.exchange.TickerPrice()
			if err != nil {
				fmt.Println("请求币安最新数据错误，10s后重试", err.Error())
				time.Sleep(10 * time.Second)
				continue
			}

			// 发送消息
			RealTimeTickerPrice = ticker
			c.Service.Send(&utils.Message{
				Type: BroadcastMessageTypeTickerPrice,
				Data: ticker,
			})
			time.Sleep(second * time.Second)
		}
	}()
}
