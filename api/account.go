package api

import (
	db "SimpleBank/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

// 声明一个创建账户请求的结构体，接收用户的请求
type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

// 为 Server 对象添加 createAccount 功能，Server 接收到用户请求，进行创建账户
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	// 将用户请求字段进行自动验证 (表单类型的参数)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 验证失败，返回 400 状态码和 JSON 格式的错误信息
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// 通过验证，则赋值给数据库创建账户的参数变量
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	// 调用 Server.store.CreateAccount 创建账户
	account, err := server.store.CreateAccount(ctx, arg)
	// 若创建账户时产生错误，则是数据库内部出错返回 500 状态码和 JSON 格式的错误信息
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 若没有产生错误，返回 200 状态码以及创建成功的账户
	ctx.JSON(http.StatusOK, account)
}

// 声明一个查找账户请求的结构体，接收用户的请求
type getAccountRequest struct {
	// 声明参数绑定，且最小值为 1
	ID int64 `uri:"id" binding:"required,min=1"`
}

// 为 Server 对象添加根据 URI 参数中的 id ，获取对应的账户的功能
func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	// 将用户请求字段进行自动验证（URI 参数类型）
	if err := ctx.ShouldBindUri(&req); err != nil {
		// 验证失败，返回 400 状态码和 JSON 格式的错误信息
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 调用 Server.store.GetAccountForUpdate 获取参数
	account, err := server.store.GetAccountForUpdate(ctx, req.ID)
	// 若获取账户时产生错误，可能是不存在该账户或者数据库内部出现错误
	if err != nil {
		// 若产生的错误为 sql.ErrNoRows ，则为账户不存在的错误，返回 404 状态码
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		// 若不是不存在该账户的错误，则是数据库内部出错返回 500 状态码和 JSON 格式的错误信息
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 若没有产生错误，返回 200 状态码以及成功查询到的账户
	ctx.JSON(http.StatusOK, account)
}

// 声明一个分页展示账户请求的结构体，接收用户的请求
type listAccountRequest struct {
	// 声明参数绑定，且最小值为 1
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// 为 Server 对象添加分页展示账户的功能
func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	// 将用户请求字段进行自动验证（查询参数类型）
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// 验证失败，返回 400 状态码和 JSON 格式的错误信息
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// 通过验证，则赋值给数据库分页展示账户的参数变量
	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	// 调用 Server.store.ListAccounts 分页展示账户
	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		// 若分页展示账户时产生错误，则是数据库内部出现错误返回 500 状态码和 JSON 格式的错误信息
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 若没有产生错误，返回 200 状态码以及分页展示的账户
	ctx.JSON(http.StatusOK, accounts)
}
