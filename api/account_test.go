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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	// 调用 randomAccount 函数创建随机账号返回给 account 变量
	account := randomAccount()

	// 创建测试用例表
	testCases := []struct {
		name      string
		accountID int64
		// 构建 stubs 的方式
		buildStubs func(store *mockdb.MockStore)
		// 检查 API 的输出
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// 构建 stubs
				// 这个 stubs 的定义可解释为：调用 GetAccountForUpdate 函数时，需要传入任何上下文和特定账户 ID 参数
				store.EXPECT().
					GetAccountForUpdate(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).            // 指定应调用此函数的次数
					Return(account, nil) // 通知 gomock 调用 GetAccountForUpdate 后返回一些特定的值
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// 检查响应
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			// 无法寻找到账户的情况的测试用例
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// 构建 stubs
				// 这个 stubs 的定义可解释为：调用 GetAccountForUpdate 函数时，需要传入任何上下文和特定账户 ID 参数
				store.EXPECT().
					GetAccountForUpdate(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).                           // 指定应调用此函数的次数
					Return(db.Account{}, pgx.ErrNoRows) // 通知 gomock 调用 GetAccountForUpdate 后返回一些特定的值
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// 检查响应
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			// 数据库内部错误的情况的测试用例
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// 构建 stubs
				// 这个 stubs 的定义可解释为：调用 GetAccountForUpdate 函数时，需要传入任何上下文和特定账户 ID 参数
				store.EXPECT().
					GetAccountForUpdate(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).                             // 指定应调用此函数的次数
					Return(db.Account{}, pgx.ErrTxClosed) // 通知 gomock 调用 GetAccountForUpdate 后返回一些特定的值
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// 检查响应
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			// 无效请求字段的情况的测试用例
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				// 构建 stubs
				// 这个 stubs 的定义可解释为：调用 GetAccountForUpdate 函数时，需要传入任何上下文和特定账户 ID 参数
				store.EXPECT().
					GetAccountForUpdate(gomock.Any(), gomock.Any()).
					Times(0) // 指定应调用此函数的次数
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// 检查响应
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	// 循环遍历测试用例
	for i := range testCases {
		// 声明变量 tc 存储当前测试用例的数据
		tc := testCases[i]

		// 将每个案例作为这个单元测试的一个单独的子测试运行
		t.Run(tc.name, func(t *testing.T) {
			// 创建一个 *gomock.Controller 对象赋值给 ctrl
			ctrl := gomock.NewController(t)
			// 延迟调用此控制器的 Finish 方法，检查是否所有预期被调用的方法都被调用
			defer ctrl.Finish()

			// 新建新 store
			store := mockdb.NewMockStore(ctrl)
			// 调用 tc.buildStubs 传入上面的 store， 创建 stubs
			tc.buildStubs(store)

			// 启动测试 HTTP 服务器并发送 GetAccountForUpdate 请求
			// 调用 NewServer 传入 mockStore 对象创建 HTTP 服务器
			server := NewServer(store)
			// 调用 httptest.NewRecorder 创建一个 ResponseRecorder 来记录 API 请求的响应
			recorder := httptest.NewRecorder()

			// 声明要调用的 API 的 URL 路径
			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			// 使用 GET 方法创建一个新的 HTTP 请求到该 URL
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// 调用 server.router.ServeHTTP 传入创建的 recorder 和 request 对象
			server.router.ServeHTTP(recorder, request)

			// 调用 tc.checkResponse 传入 t 和 recorder 进行响应检查
			tc.checkResponse(t, recorder)
		})
	}
}

// randomAccount 产生随机的账户用于测试
func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

// requireBodyMatchAccount 判断响应报文的 body 是否与传入的 account 内容相匹配
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	// 从响应报文的 body 中读取所有数据传入变量 data
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	// 声明一个 gotAccount 变量存储从响应报文 body 中获取的账户对象
	var gotAccount db.Account
	// 调用 json.Unmarshal 将数据解析后传入 gotAccount
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
