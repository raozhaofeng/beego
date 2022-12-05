package main

import (
	"fmt"
	"github.com/raozhaofeng/beego/components/blockchain"
)

func main() {

	ethereum := blockchain.NewEthereum()
	tx, isPadding := ethereum.TransactionByHash("0xfff4c8037cb1f15113f256fc50ec4d39fda1d362d90d535b0e3c9ea5c830e0ac")

	fmt.Println(isPadding)

	fmt.Println(tx.To())

	message, _ := ethereum.TransactionAsMessage(tx)
	fmt.Println(message.From())
}
