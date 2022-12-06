package main

import (
	"fmt"
	"github.com/raozhaofeng/beego/components/blockchain"
)

func main() {

	ethereum := blockchain.NewEthereum()
	tx, _ := ethereum.SetClient("https://rpc.ankr.com/eth_goerli").TransactionByHash("0x2afdef580c3c6924c4131378b18b448ba05d291854737e561832182d55e06b98")

	message, err := ethereum.TransactionAsMessage(tx)
	if err != nil {
		panic(err)
	}
	fmt.Println(message.From())
}
