package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
)

type Encryption struct {
}

func NewEncryption() *Encryption {
	return &Encryption{}
}

// DesEncrypt desc 加密
func (_Encryption *Encryption) DesEncrypt(origData, key []byte) (string, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	origData = _Encryption.PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	data := make([]byte, len(origData))
	blockMode.CryptBlocks(data, origData)
	return hex.EncodeToString(data), nil
}

// DesDecrypt des解密
func (_Encryption *Encryption) DesDecrypt(data string, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	encryptData, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	origData := make([]byte, len(encryptData))
	blockMode.CryptBlocks(origData, encryptData)
	origData = _Encryption.PKCS5UnPadding(origData)
	return origData, nil
}

// PKCS5Padding 填充方式PKCS5Padding
func (_Encryption *Encryption) PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// PKCS5UnPadding 对应 PKCS5Padding
func (_Encryption *Encryption) PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
