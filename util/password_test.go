package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	// 提供正确密码的测试情况
	// 调用 RandomString 生成长度为 6 的随机字符串作为未加密的密码
	password := RandomString(6)

	// 调用 HashPassword ，传入生成的未加密的密码，得到 hash 散列值
	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	// 调用 CheckPassword ，传入得到的未加密的密码和加密后的密码
	err = CheckPassword(password, hashedPassword1)
	require.NoError(t, err)

	// 提供错误密码的情况
	// 调用 RandomString 生成长度为 6 的随机字符串作为未加密的密码
	wrongPassword := RandomString(6)
	// 调用 CheckPassword ，传入得到的错误的未加密的密码和加密后的密码
	err = CheckPassword(wrongPassword, hashedPassword1)
	// 希望产生的错误类型为 bcrypt.ErrMismatchedHashAndPassword.Error()
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// 调用 HashPassword ，传入生成的未加密的密码，再次得到 hash 散列值
	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)

	// 相同的未加密密码加密后得到的 hash 散列值，应该要不同
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
