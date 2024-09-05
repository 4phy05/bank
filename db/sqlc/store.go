package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Store 提供所有方法单独或者在所有交易中组合执行数据库查询
type Store struct {
	*Queries
	// 为了能够创建一个数据库的事务
	connPool *pgxpool.Pool
}

// NewStore 创建一个 Store 对象
func NewStore(connPool *pgxpool.Pool) *Store {
	if connPool == nil {
		panic("cannot connect")
	}
	return &Store{
		Queries:  New(connPool),
		connPool: connPool,
	}
}

// execTx 	事务的执行一个操作并基于错误判断提交还是回滚操作
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// 开启一个新的事务
	tx, err := store.connPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	// 使用创建的事务，调用 New() 得到一个新的 *Queries 对象
	q := New(tx)
	// 使用得到的查询对象调用回调函数
	err = fn(q)
	// 如果产生了错误，进行回滚
	if err != nil {
		// 如果回滚也失败，整合产生的两个错误返回
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		// 如果回滚成功，只返回调用回调函数失败的错误
		return err
	}

	// 如果调用回调函数成功，提交事务并返回提交事务产生的错误信息
	return tx.Commit(ctx)
}

// TransferTxParams 结构体包含在两个账户之间转账所需要的所有输入参数
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	// must be positive
	Amount int64 `json:"amount"`
}

// TransferTxResult 包含交易事务的结果
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx 从一个账号到另一个账号执行一个交易，在一个事务中创建一条交易记录和账户条目并且更新账户余额
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	// 声明一个空的 TransferTxResult 的变量储存交易事务的结果
	var result TransferTxResult

	// 调用之前的 execTx() 函数去运行一个事务，在事务内进行 CURD 操作
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// 创建一条交易记录
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// 创建一条 from_account 的账户金额操作记录
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// 创建一条 to_account 的账户金额操作记录
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// 为了避免死锁，保证程序以相同的顺序获取锁，哪个 AccountID 更小就先修改哪个 账号的金额
		if arg.FromAccountID < arg.ToAccountID {
			// 调用 addMoney() 更改两个账户的余额
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			// 调用 addMoney() 更改两个账户的余额
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		// 以上过程都没有错误则返回 nil
		return nil
	})

	return result, err
}

// addMoney 一次性修改两个账号的余额
func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	// 更新 account1 的账户余额
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	// 更新 account2 的账户余额
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
