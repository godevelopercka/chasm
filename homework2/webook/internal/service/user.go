package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"webook_go/webook/internal/domain"
	"webook_go/webook/internal/repository"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error { // 调用下层 repository 的方法
	// 加密密码
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	// 然后就是，存起来
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) Edit(ctx context.Context, id int64, Nickname, Birthday, AboutMe string) (domain.User, error) {
	u, err := svc.repo.Edit(ctx, id, Nickname, Birthday, AboutMe)
	if err != nil {
		return domain.User{}, errors.New("请注册账号")
	}
	return u, nil
}

func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	u, err := svc.repo.Profile(ctx, id)
	if err != nil {
		return domain.User{}, errors.New("暂无提交个人信息")
	}
	return u, nil
}
