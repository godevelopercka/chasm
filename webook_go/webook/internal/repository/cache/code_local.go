package cache

import (
	"context"
	"sync" // 导入sync包，用于实现线程安全的并发控制
)

// CodeCache 是一个接口，定义了设置和验证验证码的方法
type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

// LocalCodeCache 是一个结构体，用于实现本地验证码缓存
type LocalCodeCache struct {
	cache map[string]map[string]string // 用于存储验证码的map，其中外层map的键为业务标识，内层map的键为手机号，值为验证码
	mutex sync.RWMutex                 // 读写锁，确保并发安全
}

// NewLocalCodeCache 是一个函数，用于创建一个新的LocalCodeCache实例
func NewLocalCodeCache() CodeCache {
	return &LocalCodeCache{
		cache: make(map[string]map[string]string),
	}
}

// Set 方法用于将验证码存储到缓存中
// 如果对应业务的验证码map不存在，则先创建它
func (c *LocalCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.cache[biz]; !ok {
		c.cache[biz] = make(map[string]string)
	}
	c.cache[biz][phone] = code
	return nil
}

// Verify 方法用于验证缓存中的验证码是否匹配
// 先通过读写锁获取读锁，检查对应业务的验证码map是否存在以及验证码是否匹配
// 如果匹配成功，则删除缓存中的验证码
func (c *LocalCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if _, ok := c.cache[biz]; !ok {
		return false, ErrUnknowForCode // 返回未知业务错误
	}
	if code, ok := c.cache[biz][phone]; ok && code == inputCode {
		delete(c.cache[biz], phone) // 验证成功后删除缓存中的验证码
		return true, nil
	}
	return false, nil // 返回验证失败错误（这里应该是自定义的错误类型ErrVerificationFailed）
}
