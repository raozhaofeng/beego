package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Random 随机
type Random struct {
	runes []rune
}

// NewRandom 随机数
func NewRandom() *Random {
	return &Random{
		runes: []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"),
	}
}

// OrderSn 随机生成订单号
func (c *Random) OrderSn() string {
	firstCode := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	nowTime := time.Now()
	return firstCode[nowTime.Year()-2020] +
		strings.ToTitle(strconv.FormatInt(int64(nowTime.Month()), 16)) +
		strconv.Itoa(nowTime.Day()) +
		strconv.FormatInt(nowTime.Unix(), 10)[5:] +
		strconv.FormatInt(nowTime.UnixNano(), 10)[2:5] +
		fmt.Sprintf("%02d", rand.Intn(100))
}

// SetLetterRunes 设置字母基数
func (c *Random) SetLetterRunes() *Random {
	c.runes = []rune("abcdefghijklmnopqrstuvwxyz")
	return c
}

// SetNumberRunes 设置数字基数
func (c *Random) SetNumberRunes() *Random {
	c.runes = []rune("0123456789")
	return c
}

// String 随机字符串
func (c *Random) String(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = c.runes[rand.Intn(len(c.runes))]
	}
	return string(b)
}
