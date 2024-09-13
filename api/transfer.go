package api

import (
	db "SimpleBank/db/sqlc"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

// 声明一个账户之间进行交易请求的结构体，接收用户的请求
type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	TOAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

// 为 Server 对象添加 createTransfer 功能，Server 接收到用户请求，账户之间进行交易
func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	// 将用户请求字段进行自动验证 (表单类型的参数)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 验证失败，返回 400 状态码和 JSON 格式的错误信息
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 调用 server.validAccount ，检验指定 FromAccountID 和 TOAccountID 的账户是否存在，以及货币类型是否对应
	if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}
	if !server.validAccount(ctx, req.TOAccountID, req.Currency) {
		return
	}

	// 通过验证，则赋值给数据库创建账户之间交易的参数变量
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.TOAccountID,
		Amount:        req.Amount,
	}

	// 调用 Server.store.TransferTx 进行账户之间的交易
	result, err := server.store.TransferTx(ctx, arg)
	// 若账户之间进行交易时产生错误，则是数据库内部出错返回 500 状态码和 JSON 格式的错误信息
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 若没有产生错误，返回 200 状态码以及成功交易的结果
	ctx.JSON(http.StatusOK, result)
}

// validAccount 检验指定 accountID 的账户是否存在，以及货币类型是否对应
func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccountForUpdate(ctx, accountID)
	if err != nil {
		// 若是未查找到账户的错误，返回 404 状态码和 JSON 格式的错误信息
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		// 否则为数据库内部的错误，返回 500 状态码和 JSON 格式的错误信息
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	// 若没有错误，检验账户的货币类型是否和输入一致
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		// 若账户货币类型和输入不一致，返回 400 状态码和 JSON 格式的错误信息
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	// 若没有产生任何错误，返回 true
	return true
}
