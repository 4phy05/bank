package api

import (
	mockdb "SimpleBank/db/mock"
	db "SimpleBank/db/sqlc"
	"SimpleBank/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// 构建自定义的匹配器
type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	// 将 x 转换为 db.CreateUserParams 对象
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	// 若转换成功，检查预期的未加密密码是否和参数中的加密后的密码相匹配
	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	// 若相匹配，则将预期参数的 HashedPassword 值设置为输入参数的 HashedPassword值
	e.arg.HashedPassword = arg.HashedPassword
	// 利用反射比较预期参数和输入参数
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

// 用于返回 EqCreateUserParamsMatcher 匹配器的实例
func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	// 创建一个随机用户，返回用户对象和未加密的密码
	user, password := randomUser(t)

	// 创建测试用例
	testCases := []struct {
		name string
		body gin.H
		// 构建 stubs 的方式
		buildStubs func(store *mockdb.MockStore)
		// 检查 API 的输出
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// 创建一个输入参数， gomock.Eq(arg) 要匹配 arg 这个输入参数
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		// 将每个案例作为这个单元测试的一个单独的子测试运行
		t.Run(tc.name, func(t *testing.T) {
			// 创建一个 *gomock.Controller 对象赋值给 ctrl
			ctrl := gomock.NewController(t)
			// 延迟调用此控制器的 Finish 方法，检查是否所有预期被调用的方法都被调用
			defer ctrl.Finish()

			// 新建新 store
			store := mockdb.NewMockStore(ctrl)
			// 调用 tc.buildStubs 传入上面的 store ，创建 stubs
			tc.buildStubs(store)

			// 启动测试 HTTP 服务器并发送 CreateUser 请求
			// 调用 NewServer 传入 mockStore 对象创建 HTTP 服务器
			server := NewServer(store)
			// 调用 httptest.NewRecorder 创建一个 ResponseRecorder 来记录 API 请求的响应
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

// randomUser 产生随机的用户用于测试，返回用户对象和未加密的密码
func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return
}

// requireBodyMatchUser 判断响应报文的 body 是否与传入的 user 内容相匹配
func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
