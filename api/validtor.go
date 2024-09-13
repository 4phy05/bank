package api

import (
	"SimpleBank/util"

	"github.com/go-playground/validator/v10"
)

// 声明 validCurrency 为验证器函数类型并且进行实例化
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	// 调用 fieldLevel.Field 返回一个反射值，通过调用 Interface 以空接口形式获取字段的值，并尝试转换为 string 类型
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// 检查是否支持该货币
		return util.IsSupportedCurrency(currency)
	}

	// 若 ok 为 false 表示该字段不是 string 类型，返回 false
	return false
}
