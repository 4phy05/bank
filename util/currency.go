package util

// 定义所有支持的货币类型
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

// IsSupportedCurrency 对输入的 currency 进行验证，若为支持的货币类型，返回 true ，否则返回 false
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}
