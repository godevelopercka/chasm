package web

//import (
//	"bytes"
//	"context"
//	"errors"
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"go.uber.org/mock/gomock"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//	"webook_go/webook/internal/domain"
//	"webook_go/webook/internal/service"
//	svcmocks "webook_go/webook/internal/service/mocks"
//)
//
//func TestUserHandler_SignUp(t *testing.T) {
//	testCase := []struct {
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
//				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
//					Email:    "123456@qq.com",
//					Password: "hello@world123",
//				}).Return(nil)
//				// 注册成功是 return nil
//				return usersvc
//			},
//			reqBody: `{
//   "email": "123456@qq.com",
//   "password": "hello@world123",
//	"confirmPassword": "hello@world123"
//}
//`,
//			wantCode: http.StatusOK,
//			wantBody: "注册成功",
//		},
//		{
//			name: "参数不对，bind 失败",
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				usersvc := svcmocks.NewMockUserService(ctrl)
//				// 注册成功是 return nil
//				return usersvc
//			},
//			reqBody: `
//{
//   "email": "123@qq.com",
//   "password": "hello@world123"
//`,
//			wantCode: http.StatusBadRequest,
//		},
//		{
//			name: "邮箱格式不正确",
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				usersvc := svcmocks.NewMockUserService(ctrl)
//				return usersvc
//			},
//			reqBody: `
//{
//   "email": "123@qm",
//   "password": "hello@world123",
//	"confirmPassword": "hello@world123"
//}
//`,
//			wantCode: http.StatusOK,
//			wantBody: "邮箱格式不正确",
//		},
//		{
//			name: "两次输入的密码不一致",
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				usersvc := svcmocks.NewMockUserService(ctrl)
//				return usersvc
//			},
//			reqBody: `
//{
//   "email": "123@qq.com",
//   "password": "hello@world1234",
//	"confirmPassword": "hello@world123"
//}
//`,
//			wantCode: http.StatusOK,
//			wantBody: "两次输入的密码不一致",
//		},
//		{
//			name: "密码格式不对",
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				usersvc := svcmocks.NewMockUserService(ctrl)
//				return usersvc
//			},
//			reqBody: `
//{
//   "email": "123@qq.com",
//   "password": "hello123",
//	"confirmPassword": "hello123"
//}
//`,
//			wantCode: http.StatusOK,
//			wantBody: "密码必须大于8位，包含数字、英文字母、特殊字符",
//		},
//		{
//			name: "邮箱冲突",
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				usersvc := svcmocks.NewMockUserService(ctrl)
//				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
//					Email:    "123@qq.com",
//					Password: "hello@world123",
//				}).Return(service.ErrUserDuplicateEmail)
//				// 注册成功是 return nil
//				return usersvc
//			},
//			reqBody: `
//{
//   "email": "123@qq.com",
//   "password": "hello@world123",
//	"confirmPassword": "hello@world123"
//}
//`,
//			wantCode: http.StatusOK,
//			wantBody: "邮箱冲突",
//		},
//		{
//			name: "系统异常",
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				usersvc := svcmocks.NewMockUserService(ctrl)
//				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
//					Email:    "123@qq.com",
//					Password: "hello@world123",
//				}).Return(errors.New("随便一个 error"))
//				// 注册成功是 return nil
//				return usersvc
//			},
//			reqBody: `
//{
//   "email": "123@qq.com",
//   "password": "hello@world123",
//	"confirmPassword": "hello@world123"
//}
//`,
//			wantCode: http.StatusOK,
//			wantBody: "系统异常",
//		},
//	}
//
//	for _, tc := range testCase {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//			server := gin.Default()
//			// 用不上 codeSvc
//			h := NewUserHandler(tc.mock(ctrl), nil, nil)
//			h.RegisterRoutes(server)
//
//			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
//			require.NoError(t, err)
//			// 数据是 JSON 格式
//			req.Header.Set("Content-Type", "application/json")
//			// 这里你就可以继续使用 req
//
//			resp := httptest.NewRecorder()
//			t.Log(resp)
//			// 这就是 HTTP 请求进去 GIN 框架的入口
//			// 当你这样调用的时候，GIN 就会处理这个请求
//			// 响应写回到 resp 里
//			server.ServeHTTP(resp, req)
//
//			assert.Equal(t, tc.wantCode, resp.Code)
//			assert.Equal(t, tc.wantBody, resp.Body.String())
//		})
//	}
//}
//
//func TestMock(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	usersvc := svcmocks.NewMockUserService(ctrl)
//
//	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
//
//	err := usersvc.SignUp(context.Background(), domain.User{
//		Email: "123@qq.com",
//	})
//	t.Log(err)
//}
//
//var (
//	server *gin.Engine
//)
//
//func init() {
//	// 在测试用例之外创建 Gin Engine 实例
//	server = gin.Default()
//}
//
//func TestUserHandler_LoginSMS(t *testing.T) {
//	testCases := []struct {
//		name     string
//		mock     func(ctrl *gomock.Controller) service.CodeService
//		reqBody  string
//		wantCode int
//		wantBody Result
//	}{
//		{
//			name: "验证码校验通过",
//			mock: func(ctrl *gomock.Controller) service.CodeService {
//				codeSvc := svcmocks.NewMockCodeService(ctrl)
//				codeSvc.EXPECT().Verify(context.Background(), "phone_code:18842003200:123456", "18842003200", "123456").Return(true, nil)
//				return codeSvc
//			},
//			reqBody:  `{"ctx": context.Background(), biz": "phone_code:15620203030:12345", "phone": "15620203030", "inputCode": "123456"}`,
//			wantCode: 200,
//			wantBody: Result{
//				Code: 4,
//				Msg:  "验证码校验通过",
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//			//server := gin.Default()
//			h := NewUserHandler(nil, tc.mock(ctrl))
//			h.RegisterRoutes(server)
//			req, err := http.NewRequest(http.MethodPost, "/users/login_sms", bytes.NewBuffer([]byte(tc.reqBody)))
//			require.NoError(t, err)
//			req.Header.Set("Content-Type", "application/json")
//			resp := httptest.NewRecorder()
//			server.ServeHTTP(resp, req)
//			assert.Equal(t, tc.wantCode, resp.Code)
//			assert.JSONEq(t, tc.wantBody.Msg, resp.Body.String())
//		})
//	}
//}
