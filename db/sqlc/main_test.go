package db

import (
	"context"
	"log"
	"os"
	"testing"

	// 用于连接 Postgres 数据库
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	// 数据库的源地址
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

// 实现了 sqlc 自动生成的方法的结构体，由实现了 DBTX 接口的变量通过 New() 产生
var testQueries *Queries

// 声明一个全局的 Store 类型变量
var testStore *Store

func TestMain(m *testing.M) {
	// 使用 pgx 库的 Connect 方法连接数据库（传入上下文和数据库地址），返回的 *pgx.Conn 属于 DBTX 接口类型
	connPool, err := pgxpool.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// 由实现了 DBTX 接口的变量通过 New() 生成实现了 sqlc 自动生成的方法的结构体
	testQueries = New(connPool)
	// 根据 *pgx.Conn 类型变量生成一个 Store 类型对象
	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
