package main

import (
	"fmt"
	"github.com/raozhaofeng/beego/components/blockchain"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {

	req, err := http.NewRequest("GET", "https://api3.binance.com/api/v3/ticker/24hr?symbol=BTCUSDT", nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))

	return

	ethereum := blockchain.NewEthereum()
	tx, _ := ethereum.SetClient("https://rpc.ankr.com/eth_goerli").TransactionByHash("0x2afdef580c3c6924c4131378b18b448ba05d291854737e561832182d55e06b98")

	message, err := ethereum.TransactionAsMessage(tx)
	if err != nil {
		panic(err)
	}
	fmt.Println(message.From())
}
