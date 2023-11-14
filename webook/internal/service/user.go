package service

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"practice/webook/internal/domain"
	"practice/webook/internal/repository"
	"practice/webook/pkg/logger"
)

var ErrUserDuplicate = repository.ErrUserDuplicate

var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
}

type SMSUserService struct {
	repo  repository.UserRepository
	redis *redis.Client
	l     logger.LoggerV1
}

func NewUserService(repo repository.UserRepository, l logger.LoggerV1) UserService {
	return &SMSUserService{
		repo: repo,
		l:    l,
	}
}

func NewUserServiceV1(repo repository.UserRepository, l *zap.Logger) UserService { // 保持依赖注入，但又没有完全注入
	return &SMSUserService{
		repo: repo,
		// 预留了变化空间
		//logger: zap.L(),
	}
}

func (svc *SMSUserService) SignUp(ctx context.Context, u domain.User) error { // context.Context 保持链路和超时控制 , 不知道返回啥就返回一个 error 这两个一定要加上
	// 要考虑加密放在哪里的问题
	// 然后就是存起来
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *SMSUserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	// 先找用户
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 比较密码了
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// 打日志 DEBUG
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *SMSUserService) FindOrCreate(ctx context.Context,
	phone string) (domain.User, error) {
	// 这时候，这个地方要怎么办
	// 这个叫做快路径
	u, err := svc.repo.FindByPhone(ctx, phone)
	// 要判断，有没有这个用户
	if err != repository.ErrUserNotFound {
		// 绝大部分请求进来这里
		// nil 会进来这里
		// 不为 ErrUserNotFound 的也会进来这里
		return u, err
	}
	// 这里，把 phone 脱敏之后打出来
	//zap.L().Info("用户未注册", zap.String("phone", phone))
	//svc.logger.Info("用户未注册", zap.String("phone", phone))
	svc.l.Info("用户未注册", logger.String("phone", phone))
	//loggerxx.Logger.Info("用户未注册", zap.String("phone", phone))
	// 在系统资源不足，触发降级后，不执行慢路径了
	//if ctx.Value("降级") == "true" {
	//	return domain.User{}, errors.New("系统降级了")
	//}
	// 这个叫做慢路径
	// 你明确知道，没有这个用户
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	// 因为这里会遇到主从延迟的问题
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *SMSUserService) FindOrCreateByWechat(ctx context.Context,
	info domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, info.OpenID)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	u = domain.User{
		WechatInfo: info,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	// 因为这里会遇到主从延迟的问题
	return svc.repo.FindByWechat(ctx, info.OpenID)
}

func (svc *SMSUserService) Profile(ctx context.Context,
	id int64) (domain.User, error) {
	u, err := svc.repo.FindById(ctx, id)
	return u, err
}

func PathsDownGrade(ctx context.Context, quick, slow func()) {
	quick()
	if ctx.Value("降级") == true {
		return
	}
	slow()
}
