package main

import (
	"fmt"
	"github.com/raozhaofeng/beego/utils"
)

func main() {
	//	加密
	data, err := utils.NewEncryption().DesEncrypt([]byte("1234512345123451234512345123451234512345"), []byte("abc11111"))
	if err != nil {
		panic(err)
	}

	fmt.Println(data)
	oldData, _ := utils.NewEncryption().DesDecrypt(data, []byte("abc11111"))

	fmt.Println(string(oldData))
}
