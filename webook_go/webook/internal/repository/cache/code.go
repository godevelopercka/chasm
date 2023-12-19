package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnknowForCode          = errors.New("我也不知发生什么了，反正是跟 code 有关")
)

// 编译器会在编译的时候，把 set_code 的代码放进来这个 luaSetCode 变量里
// 将 set_code.lua 文件嵌入到可执行文件中，并将其作为字符串赋值给 luaSetCode
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeRedisCache interface {
	Set(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeRedisCache {
	return &RedisCodeCache{
		client: client,
	}
}

// Set 方法用于在 Redis 中为指定的业务(biz)和手机号(phone)设置验证码(code)
func (c *RedisCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	// 使用 Redis 的 Eval 方法执行 Lua 脚本，设置验证码。
	// luaSetCode 是预定义的 Lua 脚本，用于在 Redis 中设置键值对。
	// c.key(biz, phone) 是生成 Redis 键的方法，键由业务标识和手机号组成。
	// code 是要设置的验证码。
	// Eval 方法执行后返回一个结果 res 和一个错误 err。
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return nil
	}
	// 根据 Lua 脚本返回的结果判断设置验证码的状态。
	switch res {
	case 0:
		// 毫无问题
		return nil
	case -1:
		// 发送太频繁
		return ErrCodeSendTooMany
	default:
		// 系统错误
		return errors.New("系统错误")
	}
}

// Verify 方法用于验证用户输入的验证码是否与 Redis 中存储的验证码匹配。
// ctx 是上下文，用于控制请求的生命周期。
// biz 是业务标识。
// phone 是手机号。
// inputCode 是用户输入的验证码。
// 如果验证码匹配，返回 true 和 nil 错误；否则返回 false 和相应的错误信息。
func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	// 使用 Redis 的 Eval 方法执行 Lua 脚本，验证验证码。
	// luaVerifyCode 是预定义的 Lua 脚本，用于验证验证码。
	// c.key(biz, phone) 是生成 Redis 键的方法，键由业务标识和手机号组成。
	// inputCode 是用户输入的验证码。
	// Eval 方法执行后返回一个结果 res 和一个错误 err。
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	// 如果执行 Lua 脚本时发生错误，则返回 false 和 err。
	if err != nil {
		return false, err
	}
	// 根据 Lua 脚本返回的结果判断验证码是否匹配。
	switch res {
	// 验证码匹配，返回 true 和 nil 错误。
	case 0:
		return true, nil
	// 验证码验证次数超过限制，返回 false 和 ErrCodeVerifyTooManyTimes 错误。
	// 正常来说，如果频繁出现这个错误，你就要告警，因为有人搞你
	case -1:
		return false, ErrCodeVerifyTooManyTimes
	// 验证码错误，返回 false 和 nil 错误。
	case -2:
		return false, nil
	// 其他未知错误，返回 false 和 ErrUnknowForCode 错误
	default:
		return false, ErrUnknowForCode
	}
}

func (c *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
