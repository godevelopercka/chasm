package web

//
//import (
//	"bytes"
//	"context"
//	"errors"
//	"github.com/gin-gonic/gin"
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"net/http"
//	"net/http/httptest"
//	"practice/webook/internal/domain"
//	"practice/webook/internal/service"
//	"testing"
//)
//
//func TestUserHandler_SignUp(t *testing.T) {
//	testCases := []struct {
//		name     string
//		mock     func(ctrl *gomock.Controller) service.UserService
//		reqBody  string
//		wantCode int
//		wantBody string
//	}{
//		{
//			name: "注册成功",
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				usersvc := svcmocks.NewMockUserService(ctrl)
//				usersvc.EXPECT().SignUp(context.Background(), domain.User{
//					Email:    "123@qq.com",
//					Password: "helloword#123",
//				}).Return(nil)
//				// 注册成功是 return nil
//				return usersvc
//			},
//			reqBody: `
//{
//	"email": "123@qq.com",
//	"password": "helloword#123",
//	"confirmPassword": "helloword#123"
//}
//`,
//			wantCode: http.StatusOK,
//			wantBody: "注册成功",
//		},
//		{
//			name: "参数不对, bind 失败",
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				usersvc := svcmocks.NewMockUserService(ctrl)
//				// 注册成功是 return nil
//				return usersvc
//			},
//			reqBody: `
//{
//	"email": "123@qq.com",
//	"password": "helloword#123",
//}
//`,
//			wantCode: http.StatusBadRequest,
//		},
//		{
//			name: "邮箱格式不对",
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				usersvc := svcmocks.NewMockUserService(ctrl)
//				// 注册成功是 return nil
//				return usersvc
//			},
//			reqBody: `
//{
//	"email": "123@q",
//	"password": "helloword#123",
//	"confirmPassword": "helloword#123"
//}
//`,
//			wantCode: http.StatusOK,
//			wantBody: "你的邮箱格式不对",
//		},
//		{
//			name: "两次输入密码不匹配",
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				usersvc := svcmocks.NewMockUserService(ctrl)
//				// 注册成功是 return nil
//				return usersvc
//			},
//			reqBody: `
//{
//	"email": "123@qq.com",
//	"password": "helloword#1234",
//	"confirmPassword": "helloword#123"
//}
//`,
//			wantCode: http.StatusOK,
//			wantBody: "两次输入的密码不一致",
//		},
//		{
//			name: "密码格式不对",
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				usersvc := svcmocks.NewMockUserService(ctrl)
//				// 注册成功是 return nil
//				return usersvc
//			},
//			reqBody: `
//{
//	"email": "123@qq.com",
//	"password": "helloword123",
//	"confirmPassword": "helloword123"
//}
//`,
//			wantCode: http.StatusOK,
//			wantBody: "密码必须大于8位，包含数字、特殊字符",
//		},
//		{
//			name: "邮箱冲突",
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				usersvc := svcmocks.NewMockUserService(ctrl)
//				usersvc.EXPECT().SignUp(context.Background(), domain.User{
//					Email:    "123@qq.com",
//					Password: "helloword#123",
//				}).Return(service.ErrUserDuplicate)
//				// 注册成功是 return nil
//				return usersvc
//			},
//			reqBody: `
//{
//	"email": "123@qq.com",
//	"password": "helloword#123",
//	"confirmPassword": "helloword#123"
//}
//`,
//			wantCode: http.StatusOK,
//			wantBody: "邮箱冲突",
//		},
//		{
//			name: "系统异常",
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				usersvc := svcmocks.NewMockUserService(ctrl)
//				usersvc.EXPECT().SignUp(context.Background(), domain.User{
//					Email:    "123@qq.com",
//					Password: "helloword#123",
//				}).Return(errors.New("随便一个 error"))
//				// 注册成功是 return nil
//				return usersvc
//			},
//			reqBody: `
//{
//	"email": "123@qq.com",
//	"password": "helloword#123",
//	"confirmPassword": "helloword#123"
//}
//`,
//			wantCode: http.StatusOK,
//			wantBody: "系统异常",
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//			server := gin.Default()
//			// 用不上 codeSvc
//			h := NewUserHandler(tc.mock(ctrl), nil)
//			h.RegisterUserRoutes(server)
//			req, err := http.NewRequest(http.MethodPost,
//				"/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
//			require.NoError(t, err)
//			// 数据是 JSON 格式
//			req.Header.Set("Content-Type", "application/json")
//			// 这里你就可以继续使用 req
//
//			resp := httptest.NewRecorder()
//			t.Log(resp)
//
//			// 这就是 HTTP 请求进去 GIN 框架的入口
//			// 当你这样调用的时候，GIN 就会处理这个请求
//			// 响应写回到 resp 里
//			server.ServeHTTP(resp, req)
//
//			assert.Equal(t, tc.wantCode, resp.Code)
//			assert.Equal(t, tc.wantBody, resp.Body.String())
//
//		})
//	}
//}
//
//func TestMock(t *testing.T) {
//	ctrl := gomock.NewController(t) // 初始化控制器
//	defer ctrl.Finish()             // 代表整个 mock 过程已经结束
//
//	usersvc := svcmocks.NewMockUserService(ctrl) //让 ctrl 控制 usersvc
//	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).
//		Return(errors.New("mock error")) // 预期发起什么调用
//	err := usersvc.SignUp(context.Background(), domain.User{
//		Email: "123@qq.com",
//	}) // context.Background 创建一个空白的上下文
//	t.Log(err)
//}
