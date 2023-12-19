package cache

import (
	"context"
	"sync" // 导入sync包，用于实现线程安全的并发控制
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

// 定义一个 LocalCodeCache 结构体，包含一个用于存储数据的 map 和一个读写锁
type LocalCodeCache struct {
	cache map[string]map[string]string
	mutex sync.RWMutex // mutex字段是一个读写锁，用于确保并发安全
}

func NewLocalCodeCache() CodeCache {
	return &LocalCodeCache{
		cache: make(map[string]map[string]string), // 初始化data字段为一个空的map
	}
}

// 定义一个方法，用于将验证码存储到缓存中
func (c *LocalCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.cache[biz]; !ok {
		c.cache[biz] = make(map[string]string)
	}
	c.cache[biz][phone] = code
	return nil
}

// 定义一个方法，用于验证缓存中的验证码是否匹配
func (c *LocalCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if _, ok := c.cache[biz]; !ok {
		return false, ErrUnknowForCode
	}
	if code, ok := c.cache[biz][phone]; ok && code == inputCode {
		delete(c.cache[biz], phone) // 验证成功后删除缓存中的验证码
		return true, nil
	}
	return false, nil
}
