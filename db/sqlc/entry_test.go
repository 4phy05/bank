package db

import (
	"SimpleBank/util"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	// 指定将要创建的条目中的项目的值
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	// 根据配置创建记录
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	// 调用 testify 包中的子包 require 的 NoError() ，判断是否没有产生错误
	require.NoError(t, err)
	// 调用 testify 包中的子包 require 的 NotEmpty() ，判断得到的 entry 是否为非空
	require.NotEmpty(t, entry)
	// 调用 testify 包中的子包 require 的 Equal() ，判断得到的 entry 中的 AccountID、Amount 是否与输入相同
	require.Equal(t, entry.AccountID, arg.AccountID)
	require.Equal(t, entry.Amount, arg.Amount)
	// 调用 testify 包中的子包 require 的 NotZero() ，判断得到的 entry 中的 ID、CreatedAt 是否为非空
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	// 创建一个 account ，用这个 account 的信息来创建 entry
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	// 创建一个 account ，用这个 account 的信息来创建 entry1
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)

	// 根据创建的 entry1 的 ID 进行查询
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	// 调用 testify 包中的子包 require 的 NoError() ，判断是否没有产生错误
	require.NoError(t, err)

	// 调用 testify 包中的子包 require 的 NotEmpty() ，判断查询得到的 entry2 是否为空
	require.NotEmpty(t, entry2)
	// 调用 testify 包中的子包 require 的 Equal() ，判断创建的 entry1 和查询得到的 entry2 是否相等
	require.Equal(t, entry1, entry2)
}

func TestListEntries(t *testing.T) {
	// 创建 1 个账号，根据这 1 个账号创建多个 entry （这里为 10 个）
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	// 指定查询参数
	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	// 根据指定的查询参数进行查询
	entries, err := testQueries.ListEntries(context.Background(), arg)
	// 调用 testify 包中的子包 require 的 NoError() ，判断是否没有产生错误
	require.NoError(t, err)
	// 调用 testify 包中的子包 require 的 Len() ，判断返回的记录数是否为指定的记录数
	require.Len(t, entries, 5)

	// 循环调用 testify 包中的子包 require 的 NotEmpty() ，判断每条记录是否不为空
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}

}
