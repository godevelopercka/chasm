package service

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/atomic"
	"math/rand"
	"practice/webook/internal/repository"
	"practice/webook/internal/service/sms"
)

var codeTplId atomic.String = atomic.String{}

var (
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
)

type CodeService interface {
	Send(ctx context.Context,
		// 区别业务场景
		biz string,
		// 这个码，谁来管？谁来生成？
		phone string) error
	Verify(ctx context.Context, biz string,
		phone string, inputCode string) (bool, error)
}

type SMSCodeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service

	//tplId string
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	codeTplId.Store("123")
	viper.OnConfigChange(func(in fsnotify.Event) {
		codeTplId.Store(viper.GetString("code.tpl.id"))
	})

	return &SMSCodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

// Send 发验证码，我需要什么参数？
func (svc *SMSCodeService) Send(ctx context.Context,
	// 区别业务场景
	biz string,
	// 这个码，谁来管？谁来生成？
	phone string) error {
	//code := "1234"
	//setToRedis(code, key, time.Minute)
	// 生成一个验证码
	code := svc.generateCode()
	// 塞进去 Redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		// 有问题
		return err
	}
	// 发送出去
	svc.smsSvc.Send(ctx, codeTplId.Load(), []string{code}, phone)
	if err != nil {
		err = fmt.Errorf("发送短信出现异常 %w", err)
	}
	//if err != nil {
	// 这个地方怎么办？
	// 这意味着，Redis 有这个验证码，但是不好意思，你没发成功，用户根本收不到
	// 我能不能删掉这个验证码？
	// 你这个 err 可能是超时的 err，你都不知道，发出了没
	// 在这里重试
	// 要重试的话，初始化的时候，传入一个自己就会重试的 smsSvc
	//}
	return err
}

// Verify 验证 验证码
func (svc *SMSCodeService) Verify(ctx context.Context, biz string,
	phone string, inputCode string) (bool, error) {
	// phone_code:login:152xxxxx
	// code:login:152xxxx
	// $biz:code:152xxxxxxx
	// user:login:code:152xxxxxx
	panic("implement me")
}

// 生成一个验证码
func (svc *SMSCodeService) generateCode() string {
	// 六位数，num 在 0, 999999 之间，包含 0 和 999999
	num := rand.Intn(1000000)
	// 不够六位的，加上前导 0
	// 000001
	return fmt.Sprintf("%6d", num)
}
