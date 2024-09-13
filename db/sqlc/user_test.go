package db

import (
	"SimpleBank/util"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	// 调用 util.RandomString 产生长度为 6 的随机未加密密码
	hashPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	// 使用 util 包中编写的随机函数生成测试用例
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	// 调用 sqlc 生成的 CreateUser 方法，进行测试
	user, err := testQueries.CreateUser(context.Background(), arg)
	// 调用 testify 包中的子包 require 的 NoError() ，判断是否没有产生错误
	require.NoError(t, err)
	// 调用 testify 包中的子包 require 的 NotEmpty() ，判断得到的 user 是否不为空
	require.NotEmpty(t, user)
	// 调用 testify 包中的子包 require 的 Equal() ，判断得到的 user 中的 Username、HashedPassword、FullName、Email 是否与输入相同
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	// 调用 testify 包中的子包 require 的 NotZero() ，判断 PasswordChangedAt、CreatedAt 是否由 Postgres 自动生成
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	// 创建随机用户并且返回创建的 user 对象
	user1 := createRandomUser(t)
	// 根据上面创建的对象的 Username 查询出对应的记录
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	// 调用 testify 包中的子包 require 的 NoError() ，判断是否没有产生错误
	require.NoError(t, err)
	// 调用 testify 包中的子包 require 的 NotEmpty() ，判断查询得到的 user2 是否为空
	require.NotEmpty(t, user2)

	// 调用 testify 包中的子包 require 的 Equal() ，判断创建的 user1 和查询得到的 user2 各项数据是否相同
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	// 调用 testify 包中的子包 require 的 WithinDuration() ，判断 PasswordChangedAt、CreatedAt 是否在允许的偏差范围内，这里为 time.Second
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

}
