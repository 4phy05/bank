package db

import (
	"SimpleBank/util"
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	// 先创建一个 user
	user := createRandomUser(t)
	// 使用 util 包中编写的随机函数和 user.Username 生成测试用例
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	// 调用 sqlc 生成的 CreateAccount 方法，进行测试
	account, err := testQueries.CreateAccount(context.Background(), arg)
	// 调用 testify 包中的子包 require 的 NoError() ，判断是否没有产生错误
	require.NoError(t, err)
	// 调用 testify 包中的子包 require 的 NotEmpty() ，判断得到的 account 是否不为空
	require.NotEmpty(t, account)
	// 调用 testify 包中的子包 require 的 Equal() ，判断得到的 account 中的 Owner、Balance、Currency 是否与输入相同
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	// 调用 testify 包中的子包 require 的 NotZero() ，判断 ID、CreatedAt 是否由 Postgres 自动生成
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccountForUpdate(t *testing.T) {
	// 创建随机用户并且返回创建的 account 对象
	account1 := createRandomAccount(t)
	// 根据上面创建的对象的 ID 查询出对应的记录
	account2, err := testQueries.GetAccountForUpdate(context.Background(), account1.ID)

	// 调用 testify 包中的子包 require 的 NoError() ，判断是否没有产生错误
	require.NoError(t, err)
	// 调用 testify 包中的子包 require 的 NotEmpty() ，判断查询得到的 account2 是否为空
	require.NotEmpty(t, account2)

	// 调用 testify 包中的子包 require 的 Equal() ，判断创建的 account1 和查询得到的 account2 各项数据是否相同
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	// 调用 testify 包中的子包 require 的 WithinDuration() ，判断 CreatedAt 是否在允许的偏差范围内，这里为 time.Second
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	// 创建随机用户并且返回创建的 account 对象
	account1 := createRandomAccount(t)

	// 使用 util 包中编写的随机函数生成测试用例
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	// 根据新创建的 account1 修改其 Balance 项
	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	// 调用 testify 包中的子包 require 的 NoError() ，判断是否没有产生错误
	require.NoError(t, err)
	// 调用 testify 包中的子包 require 的 NotEmpty() ，判断 account2 是否不为空
	require.NotEmpty(t, account2)
	// 调用 testify 包中的子包 require 的 Equal() ，判断创建的 account1 和修改后得到的 account2 各项数据是否相同（除了Balance）
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	// 调用 testify 包中的子包 require 的 Equal() ，判断修改之后 account2 中的 Balance 项是否和输入的 arg 中的相同
	require.Equal(t, account2.Balance, arg.Balance)
}

func TestDeleteAccount(t *testing.T) {
	// 创建随机用户并且返回创建的 account 对象
	account1 := createRandomAccount(t)

	// 根据创建的用户 account1 的 ID 项删除记录
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	// 调用 testify 包中的子包 require 的 NoError() ，判断是否没有产生错误
	require.NoError(t, err)

	// 根据创建的用户 account1 的 ID 项查询记录
	account2, err := testQueries.GetAccountForUpdate(context.Background(), account1.ID)
	// 调用 testify 包中的子包 require 的 Error() ，判断是否产生错误
	require.Error(t, err)
	// 调用 testify 包中的子包 require 的 EqualError() ，判断产生的错误是否为 pgx.ErrNoRows.Error()
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	// 调用 testify 包中的子包 require 的 Empty() ，判断查询到的 account2 是否为空
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	// 创建多个账号（这里指定为 10 个）
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	// 指定查询参数
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	// 根据指定的查询参数进行查询
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	// 调用 testify 包中的子包 require 的 NoError() ，判断是否没有产生错误
	require.NoError(t, err)
	// 调用 testify 包中的子包 require 的 Len() ，判断返回的记录数是否为指定的记录数
	require.Len(t, accounts, 5)

	// 循环调用 testify 包中的子包 require 的 NotEmpty() ，判断每条记录是否不为空
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
