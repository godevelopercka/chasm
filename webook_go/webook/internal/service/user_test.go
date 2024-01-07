package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
	"webook_go/webook/internal/domain"
	"webook_go/webook/internal/repository"
	repomocks "webook_go/webook/internal/repository/mocks"
)

func TestUserHandler_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		// 输入
		ctx      context.Context
		email    string
		password string

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功", // 用户名和密码是对的
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					Password: "$2a$10$QKwcZ5m2NeL36NHBdu8PR.wthInx4AFK2KTweRkKdtcBBL",
					Phone:    "15245311230",
					Ctime:    now,
				}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world123",

			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$QKwcZ5m2NeL36NHBdu8PR.wthInx4AFK2KTweRkKdtcBBL",
				Phone:    "15245311230",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world123",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "DB错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, errors.New("mock db 错误"))
				return repo
			},
			email:    "123@qq.com",
			password: "helloadfa#world123",

			wantUser: domain.User{},
			wantErr:  errors.New("mock db 错误"),
		},
		{
			name: "密码不对",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					Password: "$2a$10$QKwcZ5m2NeL36NHBdu8PR.wthInx4AFK2KTweRkKdtcBBL",
					Phone:    "15245311230",
					Ctime:    now,
				}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world123",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl), nil)
			u, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

// 获取加密后的密码，放到上面的测试用例中
func TestEncrypted(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("hello#world123"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}
