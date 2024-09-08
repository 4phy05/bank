package api

import (
	db "SimpleBank/db/sqlc"

	"github.com/gin-gonic/gin"
)

// Server 为银行服务提供所有的 HTTP 请求
type Server struct {
	// 允许处理来自客户端的 API 请求时，与数据库进行交互
	store db.Store
	// 帮助将每个 API 请求发送到正确的处理程序进行处理
	router *gin.Engine
}

// NewServer 创建一个服务器，并在服务器上设置路由
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// 为 router 添加路由处理
	// 创建账户
	router.POST("/accounts", server.createAccount)
	// 根据 ID 访问指定的账户
	router.GET("/accounts/:id", server.getAccount)
	// 分页展示账户
	router.GET("/accounts", server.listAccount)

	// 将配置好的 router 配置到 Server 上
	server.router = router
	return server
}

// 在指定的 address 上运行 HTTP 服务器，开始监听 API 请求
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse 将错误信息格式转换为 gin.H
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
