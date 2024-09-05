package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefjhijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt 生成一个随机数处在 min 和 max 之间
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString 生成一个随机字符串，长度为 n
func RandomString(n int) string {
	// 声明一个变量 sb 为 strings.Builder 内置切片类型可以写入字符
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		// 随机得到 alphabet 字母表中的一个字符
		c := alphabet[rand.Intn(k)]
		// 写入 sb 变量中
		sb.WriteByte(c)
	}

	// 返回以 string 类型返回 sb 变量
	return sb.String()
}

// RandomOwner 生成 6 个字符的随机用户名
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney 生成 0-1000 之间随机金额的钱
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency 生成 EUR/USD/CAD 三种货币类型的随机一种
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)

	return currencies[rand.Intn(n)]
}
