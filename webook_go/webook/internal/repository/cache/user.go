package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"webook_go/webook/internal/domain"
)

var ErrKeyNotExist = redis.Nil

// 定义一个接口承接下面 RedisUserCache 结构体的方法，如何在 service 里面将这个接口作为字段使用
type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
}

type RedisUserCache struct {
	// 传单机 Redis 可以
	// 传 cluster 的 Redis 也可以
	client     redis.Cmdable
	expiration time.Duration
}

// A 用到了 B, B 一定是接口 => 这个是保证面向接口
// A 用到了 B, B 一定是 A 的字段 => 规避包变量、包方法，都非常缺乏扩展性
// A 用到了 B, A 绝对不初始化 B, 而是外面注入 => 保持依赖注入(DI, Dependency Injection)和依赖反转(IOC)
func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

// 如果没有数据，返回一个特定的 error
func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	// 数据不存在，err = redis.nil
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u) // 将 val 的值，放到 u 中
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
