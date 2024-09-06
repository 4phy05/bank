package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	// 创建两个新建的账户，进行交易
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	// 打印交易前每个账号的余额
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// 指定创建交易条数和每次交易的金额（这里分别是 5 和 10）
	n := 5
	amount := int64(10)

	// 声明两个 chan 类型，接收 goroutine 并发过程中的结果
	errs := make(chan error)
	results := make(chan TransferTxResult)

	// 使用两个新创建的账户进行多次交易
	for i := 0; i < n; i++ {

		ctx := context.Background()
		go func() {
			result, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			// 将所有 goroutine 中得到的数据都存入事先声明的 chan 中
			errs <- err
			results <- result
		}()
	}

	// 声明一个变量存放每次交易为第几次交易
	existed := make(map[int]bool)
	// 检验测试的结果返回的 channel errs 和 channel results 中的内容
	for i := 0; i < n; i++ {
		// 将 errs channel 中的数据依次取出
		err := <-errs
		// 调用 testify 包中的子包 require 的 NoError() ，判断是否没有产生错误
		require.NoError(t, err)

		// 将 results channel 中的数据一次取出
		result := <-results
		// 调用 testify 包中的子包 require 的 NotEmpty() ，判断取出的内容是否为非空
		require.NotEmpty(t, result)

		// 检验 result 中的 Transfer 对象
		transfer := result.Transfer
		// 检验非空、自动填充字段是否非零和各个字段是否和填入的相等
		require.NotEmpty(t, transfer)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)

		// 根据 transfer 的 ID 去数据库查询是否创建记录成功
		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		// 判断是否没有错误产生
		require.NoError(t, err)

		// 检验 result 中的 FromEntry 对象
		fromEntry := result.FromEntry
		// 检验非空、自动填充字段是否非零和各个字段是否和填入的相等
		require.NotEmpty(t, fromEntry)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)

		// 根据 fromEntry 的 ID 去数据库查询是否创建记录成功
		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		// 判断是否没有错误产生
		require.NoError(t, err)

		// 检验 result 中的 ToEntry 对象
		toEntry := result.ToEntry
		// 检验非空、自动填充字段是否非零和各个字段是否和填入的相等
		require.NotEmpty(t, toEntry)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)

		// 根据 toEntry 的 ID 去数据库查询是否创建记录成功
		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		// 判断是否没有错误产生
		require.NoError(t, err)

		// 检验 result 中的 FromAccount 对象
		fromAccount := result.FromAccount
		// 检验非空、ID 字段是否和填入的相等
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		// 检验 result 中的 FromAccount 对象
		toAccount := result.ToAccount
		// 检验非空、ID 字段是否和填入的相等
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// 打印每次交易后每个账号的余额
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		// 检验交易过后的账户余额
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		// 检验从 account1 转出的钱是否和转入 account2 的钱相等，都为正数且交易金额为每次交易金额的整数倍
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		// 声明变量存放当前交易额为每次交易的多少倍，表示这是第几次交易
		k := int(diff1 / amount)
		// 检验此次交易是否不存在之前记录的交易中
		require.NotContains(t, existed, k)
		// 将每次为第几次交易成功存入 existed 中
		existed[k] = true
	}

	// 检验最终所有交易完成之后 account1 的账户余额
	updatedAccount1, err := testQueries.GetAccountForUpdate(context.Background(), account1.ID)
	// 检验是否没产生错误，且是否最终减少金额为 n 倍每次交易金额
	require.NoError(t, err)
	require.Equal(t, updatedAccount1.Balance+int64(n)*amount, account1.Balance)

	// 检验最终所有交易完成之后 account2 的账户余额
	updatedAccount2, err := testQueries.GetAccountForUpdate(context.Background(), account2.ID)
	// 检验是否没产生错误,且是否最终增加金额为 n 倍每次交易金额
	require.NoError(t, err)
	require.Equal(t, updatedAccount2.Balance-int64(n)*amount, account2.Balance)

	// 打印最终交易后每个账号的余额
	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	// 创建两个新建的账户，进行交易
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	// 打印交易前每个账号的余额
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// 指定创建交易条数和每次交易的金额（这里分别是 5 和 10）
	n := 10
	amount := int64(10)

	// 声明两个 chan 类型，接收 goroutine 并发过程中的结果
	errs := make(chan error)

	// 使用两个新创建的账户进行多次交易
	for i := 0; i < n; i++ {
		// 一半的交易从 account1 到 account2 ，一半的交易从 account2 到 account1
		fromAccountID := account1.ID
		toAccountID := account2.ID
		// 当交易次数位偶数时，交易从 account2 到 account1
		if i%2 == 0 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		ctx := context.Background()
		go func() {
			_, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			// 将所有 goroutine 中得到的数据都存入事先声明的 chan 中
			errs <- err
		}()
	}

	// 检验测试的结果返回的 channel errs
	for i := 0; i < n; i++ {
		// 将 errs channel 中的数据依次取出
		err := <-errs
		// 调用 testify 包中的子包 require 的 NoError() ，判断是否没有产生错误
		require.NoError(t, err)
	}

	// 检验最终所有交易完成之后 account1 的账户余额
	updatedAccount1, err := testQueries.GetAccountForUpdate(context.Background(), account1.ID)
	// 检验是否没产生错误，且是否最终余额是否没变
	require.NoError(t, err)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)

	// 检验最终所有交易完成之后 account2 的账户余额
	updatedAccount2, err := testQueries.GetAccountForUpdate(context.Background(), account2.ID)
	// 检验是否没产生错误,且是否最终余额是否没变
	require.NoError(t, err)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

	// 打印最终交易后每个账号的余额
	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
}
