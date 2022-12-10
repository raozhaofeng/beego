package exchange

import (
	"errors"
	"fmt"
	"github.com/raozhaofeng/beego/components/exchange/binance"
	"github.com/raozhaofeng/beego/utils"
)

// BroadcastService 广播服务
var BroadcastService = map[string]map[int64]*Broadcast{}

// BroadcastServiceManageService 广播服务管理
var BroadcastServiceManageService *BroadcastManage

// BroadcastManage 广播管理
type BroadcastManage struct {
}

// Broadcast 广播
type Broadcast struct {
	isRun           bool
	chListeners     chan *utils.Message
	chListenersFunc func(broadcastMessage *utils.Message) error
}

// InitBroadcastService 初始化广播服务
func InitBroadcastService() {
	// 如果币安交易所规范不存在，那么请求规范信息
	if binance.SpecificationInfo == nil {
		binance.InitSpecification()
	}

	// 如果币安交易所广播不存在，那么执行币安广播
	if binance.BroadcastService == nil {
		binance.InitBroadcastService()
		// 启动广播
		binance.BroadcastService.Service.Start()
	}

	BroadcastServiceManageService = &BroadcastManage{}
}

// StartListen 启动监听
func (c *BroadcastManage) StartListen(exchangeName string, chIndex int64) error {
	if _, ok := BroadcastService[exchangeName]; !ok {
		return errors.New("NotFoundExchangeName")
	}
	if _, ok := BroadcastService[exchangeName][chIndex]; !ok {
		return errors.New("NotFoundExchangeNameChIndex")
	}
	if !BroadcastService[exchangeName][chIndex].isRun {
		return nil
	}

	go func() {
		for broadcastMessage := range BroadcastService[exchangeName][chIndex].chListeners {
			err := BroadcastService[exchangeName][chIndex].chListenersFunc(broadcastMessage)
			if err != nil {
				fmt.Println(err)
				BroadcastService[exchangeName][chIndex].isRun = false
				_ = c.RemoveListenerTickerPrice(exchangeName, chIndex)
			}
		}
	}()
	return nil
}

// ListenerTickerPrice 最新价格监听者
func (c *BroadcastManage) ListenerTickerPrice(exchangeName string, chIndex int64, chListenersFunc func(broadcastMessage *utils.Message) error) error {
	var chListeners chan *utils.Message
	switch exchangeName {
	case BinanceExchange:
		chListeners = binance.BroadcastService.Service.Listener(chIndex)
	default:
		return errors.New("NotFoundExchange")
	}

	if _, ok := BroadcastService[exchangeName]; !ok {
		BroadcastService[exchangeName] = make(map[int64]*Broadcast, 0)
	}
	if _, ok := BroadcastService[exchangeName][chIndex]; !ok {
		BroadcastService[exchangeName][chIndex] = &Broadcast{
			chListeners:     chListeners,
			chListenersFunc: chListenersFunc,
		}

		// 直接启动监听模式
		_ = c.StartListen(exchangeName, chIndex)
	}
	return nil
}

// RemoveListenerTickerPrice 删除监听者
func (c *BroadcastManage) RemoveListenerTickerPrice(exchangeName string, chIndex int64) error {
	if _, ok := BroadcastService[exchangeName]; !ok {
		return errors.New("NotFoundExchangeName")
	}
	if _, ok := BroadcastService[exchangeName][chIndex]; !ok {
		return errors.New("NotFoundExchangeNameChIndex")
	}

	switch exchangeName {
	case BinanceExchange:
		binance.BroadcastService.Service.RemoveListener(chIndex, BroadcastService[exchangeName][chIndex].chListeners)
	default:
		return errors.New("NotFoundExchangeName")
	}
	return nil
}
