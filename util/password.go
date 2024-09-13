package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 计算输入密码的 bcrypt 散列值
func HashPassword(password string) (string, error) {
	// 传入 byte 切片类型的 password ，和 bcrypt 默认的 cost ，得到加密后的密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// 若产生错误，返回空的加密后的密码和包装后的错误（方便进行错误的维护）
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword 与提供的 hashedPassword 进行对比，检查输入的密码是否正确
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
