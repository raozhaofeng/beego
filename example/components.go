package main

import (
	"fmt"
	"github.com/raozhaofeng/beego/components/blockchain"
	"math"
	"math/big"
)

func main() {

	big.NewFloat(7.898989)

	m := new(big.Float)

	s, _ := m.Mul(big.NewFloat(7.898988984838123821839), big.NewFloat(math.Pow10(6))).Int(&big.Int{})
	fmt.Println(s.String())

	//balance := new(big.Float)
	//balance.SetFloat64(10.9)

	//s, _ := utils.NewEncryption().DesEncrypt([]byte("982365c3d87cfbd772b3f066ff89578186189970d99106b8bd48532923fdf70d"), []byte("Aa123098"))

	//fmt.Println(s)
	return
	ethereum := blockchain.NewEthereum()
	tx, _ := ethereum.SetClient("https://rpc.ankr.com/eth_goerli").TransactionByHash("0x2afdef580c3c6924c4131378b18b448ba05d291854737e561832182d55e06b98")

	message, err := ethereum.TransactionAsMessage(tx)
	if err != nil {
		panic(err)
	}
	fmt.Println(message.From())
}
