package main

import (
	"fmt"
	"github.com/raozhaofeng/beego/components/blockchain"
)

func main() {

	//fmt.Println(blockchain.NewEthereum().Balance("0x13ff0c6B6651FdaEc3EA8C424fb65791fA92Bd9a"))

	//	测试私钥
	//	0xb62cb3de57d2000184d7b4850b3441eb2801ef7b52c1c5a3480fbc0595627164

	//privateHex := "b62cb3de57d2000184d7b4850b3441eb2801ef7b52c1c5a3480fbc0595627164"

	//	usdt合约地址
	//contractAddress := "0xdac17f958d2ee523a2206206994597c13d831ec7"

	//	转账地址
	//address := "0x03a14F3A3bfE61846805cFf96118651449d84876"

	//	交易hex
	hex := "0xb0c76c132d977d9ec2e6c46115b4108873531c375daef95466f29b8eb7dbfb9e"

	//	转账数量
	//var value int64 = 1

	status, err := blockchain.NewEthereum().TransactionReceiptStatus(hex)
	if err != nil {
		panic(err)
	}

	fmt.Println(status)
}
