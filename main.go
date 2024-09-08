package main

import (
	db "SimpleBank/db/sqlc"
	"SimpleBank/util"
	"context"
	"log"

	"SimpleBank/api"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	// 在当前文件夹中读取配置文件，解析配置值到 config 变量
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// 使用 pgx 库的 Connect 方法连接数据库（传入上下文和数据库地址），返回的 *pgx.Conn 属于 DBTX 接口类型
	conn, err := pgxpool.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// 根据 *pgx.Conn 类型变量生成一个 Store 类型对象
	store := db.NewStore(conn)

	// 根据生成的 store 创建一个 sever
	server := api.NewServer(store)

	// 启动上面创建的 server，并监听指定的地址
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("connot start server:", err)
	}
}
