package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/raozhaofeng/beego/components/exchange/interfaces"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// BaseURL 基础接口
const BaseURL = "https://api.binance.com"

const (
	// TickerPriceURL 最新交易对价格列表路由
	TickerPriceURL = "/api/v3/ticker/price"

	// KlineURL k线路由
	KlineURL = "/api/v3/klines"

	// SpecificationURL 行情规范路由
	SpecificationURL = "/api/v3/exchangeInfo"

	// OrderURL 订单路由
	OrderURL = "/api/v3/order"

	// TestOrderURL 测试订单路由
	TestOrderURL = "/api/v3/order/test"
)

// httpRespError 请求返回数据错误
type httpRespError struct {
	Code int64  `json:"code"` //	错误代码
	Msg  string `json:"msg"`  //	错误消息
}

// Exchange 币安交易所
type Exchange struct {
	conf   interfaces.Conf //	配置文件
	isTest bool            //	是否测试
}

// NewExchange 新建币安交易所模型
func NewExchange() *Exchange {
	return &Exchange{}
}

// Conf 配置交易所配置
func (c *Exchange) Conf(conf interfaces.Conf) interfaces.Exchange {
	c.conf = conf
	return c
}

// TestOrder 测试订单
func (c *Exchange) TestOrder() interfaces.Exchange {
	c.isTest = true
	return c
}

// SpecificationInfo 交易所规范信息
func (c *Exchange) SpecificationInfo() (*Info, error) {
	resp, err := c.httpGet(SpecificationURL)
	if err != nil {
		return nil, err
	}

	data := new(Info)
	err = json.Unmarshal(resp, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// TickerPrice 请求交易对最新价格列表
func (c *Exchange) TickerPrice() (interfaces.Ticker, error) {
	resp, err := c.httpGet(TickerPriceURL)
	if err != nil {
		return nil, err
	}

	var data = NewTicker()
	err = json.Unmarshal(resp, &data.TickerList)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Buy 买入
func (c *Exchange) Buy(p map[string]interface{}) (interfaces.Resp, error) {
	params := map[string]interface{}{
		"side":       "BUY",
		"recvWindow": 60000,
		"timestamp":  time.Now().Unix() * 1000,
	}
	for k, v := range p {
		params[k] = v
	}
	if c.isTest {
		return c.httpError(c.httpPost(TestOrderURL, params))
	}
	return c.httpError(c.httpPost(OrderURL, params))
}

// MarketBuy 市价买入
func (c *Exchange) MarketBuy(symbol string, quoteOrderQty float64) (interfaces.Resp, error) {
	p := map[string]interface{}{
		"symbol":        symbol,
		"quoteOrderQty": quoteOrderQty,
		"type":          "MARKET",
	}
	return c.Buy(p)
}

// Sell 卖出
func (c *Exchange) Sell(p map[string]interface{}) (interfaces.Resp, error) {
	params := map[string]interface{}{
		"side":       "SELL",
		"recvWindow": 60000,
		"timestamp":  time.Now().Unix() * 1000,
	}
	for k, v := range p {
		params[k] = v
	}
	if c.isTest {
		return c.httpError(c.httpPost(TestOrderURL, params))
	}
	return c.httpError(c.httpPost(OrderURL, params))
}

// MarketSell 市价卖出
func (c *Exchange) MarketSell(symbol string, quantity float64) (interfaces.Resp, error) {
	filterQuantity, err := c.filterMinQtyQuantity(symbol, quantity)
	if err != nil {
		return nil, err
	}
	p := map[string]interface{}{
		"symbol":   symbol,
		"quantity": filterQuantity,
		"type":     "MARKET",
	}
	return c.Sell(p)
}

// httpGet GET请求
func (c *Exchange) httpGet(uri string) ([]byte, error) {
	req, err := http.NewRequest("GET", BaseURL+uri, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

// httpPost POST请求
func (c *Exchange) httpPost(uri string, p map[string]interface{}) ([]byte, error) {
	queryStr := c.mapToQueryParams(p)
	sha := c.sign(queryStr)
	queryStr += "&signature=" + sha

	req, err := http.NewRequest("POST", BaseURL+uri, strings.NewReader(queryStr))
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req.Header.Add("X-MBX-APIKEY", c.conf.AppKey)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

// httpError 判断返回是否错误
func (c *Exchange) httpError(resp []byte, err error) (interfaces.Resp, error) {
	if err != nil {
		return nil, err
	}

	respError := new(httpRespError)
	_ = json.Unmarshal(resp, respError)
	if respError.Code != 0 {
		return nil, errors.New(fmt.Sprintf("错误代码:%v, 错误信息:%v", respError.Code, respError.Msg))
	}

	respStruct := new(Resp)
	err = json.Unmarshal(resp, respStruct)
	if err != nil {
		return nil, err
	}
	return respStruct, nil
}

// sign 加密方法
func (c *Exchange) sign(queryStr string) string {
	h := hmac.New(sha256.New, []byte(c.conf.SecretKey))
	h.Write([]byte(queryStr))
	return hex.EncodeToString(h.Sum(nil))
}

// mapToQueryParams map转请求参数
func (c *Exchange) mapToQueryParams(mapList map[string]interface{}) string {
	var s []string
	for k, v := range mapList {
		s = append(s, fmt.Sprintf("%v=%v", k, v))
	}
	return strings.Join(s, "&")
}

// filterMinQtyQuantity 过滤最低数量
func (c *Exchange) filterMinQtyQuantity(symbol string, quantity float64) (float64, error) {
	if _, ok := SpecificationSymbolInfo[symbol]; !ok {
		return 0, errors.New("NotSpecification")
	}
	if _, ok := SpecificationSymbolInfo[symbol].FiltersMap["LOT_SIZE"]; !ok {
		return 0, errors.New("Not[LOT_SIZE]")
	}

	minQty := SpecificationSymbolInfo[symbol].FiltersMap["LOT_SIZE"]["minQty"].(string)
	minQtyList := strings.Split(minQty, ".")
	quantityStr := strconv.FormatFloat(quantity, 'f', -1, 64)
	quantityList := strings.Split(quantityStr, ".")

	// 如果是整数，那么直接返回当前数量不在验证
	if len(quantityList) != 2 {
		return quantity, nil
	}

	newQuantityStr := ""
	if minQtyList[0] != "0" {
		newQuantityStr = quantityList[0]
	} else {
		newQuantityStr = quantityList[0] + "."
		if len(quantityList[1]) > len(minQtyList[1]) {
			newQuantityStr += quantityList[1][:len(minQtyList[1])]
		} else {
			newQuantityStr += quantityList[1]
		}
	}

	newQuantity, _ := strconv.ParseFloat(newQuantityStr, 64)
	return newQuantity, nil
}
