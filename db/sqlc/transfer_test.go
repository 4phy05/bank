package db

import (
	"SimpleBank/util"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, from_account, to_account Account) Transfer {
	// 指定将要创建的条目中的项目的值
	arg := CreateTransferParams{
		FromAccountID: from_account.ID,
		ToAccountID:   to_account.ID,
		Amount:        util.RandomMoney(),
	}

	// 根据配置创建记录
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	// 调用 testify 包中的子包 require 的 NoError() ，判断是否没有产生错误
	require.NoError(t, err)
	// 调用 testify 包中的子包 require 的 NotEmpty() ，判断返回的记录是否为非空
	require.NotEmpty(t, transfer)

	// 调用 testify 包中的子包 require 的 Equal() ，判断返回的 transfer 的 FromAccountID 是否与 arg 相等
	require.Equal(t, transfer.FromAccountID, from_account.ID)
	// 调用 testify 包中的子包 require 的 Equal() ，判断返回的 transfer 的 ToAccountID 是否与 arg 相等
	require.Equal(t, transfer.ToAccountID, to_account.ID)
	// 调用 testify 包中的子包 require 的 Equal() ，判断返回的 transfer 的 Amount 是否与 arg 相等
	require.Equal(t, transfer.Amount, arg.Amount)

	// 调用 testify 包中的子包 require 的 NotZero() ，判断得到的 transfer 中的 ID、CreatedAt 是否为非空
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	// 创建两个账号分别充当转账账号和进账账号，通过这两个账号创建交易记录
	from_account := createRandomAccount(t)
	to_account := createRandomAccount(t)
	createRandomTransfer(t, from_account, to_account)
}

func TestGetTransfer(t *testing.T) {
	// 创建两个账号分别充当转账账号和进账账号，通过这两个账号创建交易记录
	from_account := createRandomAccount(t)
	to_account := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, from_account, to_account)

	// 根据创建的 transfer1 的 ID 进行查询
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	// 调用 testify 包中的子包 require 的 NoError() ，判断是否没有产生错误
	require.NoError(t, err)

	// 调用 testify 包中的子包 require 的 Equal() ，判断 transfer2 是否为非空
	require.NotZero(t, transfer2)
	// 调用 testify 包中的子包 require 的 Equal() ，判断 transfer1 与 transfer2 是否相等
	require.Equal(t, transfer1, transfer2)
}

func TestListTransfers(t *testing.T) {
	// 创建两个账号分别充当转账账号和进账账号，通过这两个账号创建多条交易记录（这里是 10 条）
	from_account := createRandomAccount(t)
	to_account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, from_account, to_account)
	}

	// 指定查询参数
	arg := ListTransfersParams{
		FromAccountID: from_account.ID,
		ToAccountID:   to_account.ID,
		Limit:         5,
		Offset:        5,
	}

	// 根据指定的查询参数进行查询
	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	// 调用 testify 包中的子包 require 的 NoError() ，判断是否没有产生错误
	require.NoError(t, err)

	// 调用 testify 包中的子包 require 的 Len() ，判断返回的记录数是否为指定的记录数
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		// 调用 testify 包中的子包 require 的 NotEmpty() ，判断每条记录是否不为空
		require.NotEmpty(t, transfer)
		// 调用 testify 包中的子包 require 的 Equal() ，判断 transfer 的 FromAccountID 是否与 arg 相等
		require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
		require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
	}
}
