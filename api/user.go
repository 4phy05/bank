package api

import (
	db "SimpleBank/db/sqlc"
	"SimpleBank/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
)

// 声明一个创建用户请求的结构体，接收用户的请求
type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// 声明一个创建用户响应的结构体，用于响应用户创建用户的请求
type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

// 为 Server 对象添加 createUser 功能，Server 接收到用户请求，进行创建账户
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	// 将用户请求字段进行自动验证 (表单类型的参数)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 验证失败，返回 400 状态码和 JSON 格式的错误信息
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 通过验证，则调用 util.HashPassword 生成加密后的密码
	hashedPassword, err := util.HashPassword(req.Password)
	// 若加密过程产生错误，则返回 500 状态码以及 JSON 格式的错误信息
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// 赋值给数据库创建账户的参数变量
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	// 调用 Server.store.CreateUser 创建账户
	user, err := server.store.CreateUser(ctx, arg)
	// 若创建账户时产生错误，则是可能是数据库内部出错或者违反约束
	if err != nil {
		// 若出现错误，尝试将错误转换为 *pgconn.PgError 类型
		if pqErr, ok := err.(*pgconn.PgError); ok {
			// 根据返回的状态码，返回 403 状态码和 JSON 格式的错误信息
			switch pqErr.Code {
			case "23505":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		// 数据库内部出错，返回 500 状态码和 JSON 格式的错误信息
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
	// 若没有产生错误，返回 200 状态码以及成功创建账户的响应
	ctx.JSON(http.StatusOK, rsp)
}
