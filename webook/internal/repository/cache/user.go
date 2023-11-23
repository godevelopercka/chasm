package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"practice/webook/internal/domain"
	"time"
)

// 最差的方法
//type Cache interface {
//	GetUser(ctx context.Context, id int64) (domain.User, error)
//	// 读取文章
//	GetArticle(ctx context.Context, aid int64)
//	// 还有别的业务
//	// ...
//}

// 最好的方法
//type CacheV1 interface {
//	// 你的中间件团队去做的
//	Get(ctx context.Context, key string) (any, error)
//}
//
//type RedisUserCache struct {
//	cache CacheV1
//}
//
//func (u *RedisUserCache) GetUser(ctx context.Context, id int64) (domain.User, error) {
//
//}

// 序列化与反序列化

var ErrKeyNotExist = redis.Nil

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

// 只要 error 为 nil, 就认为缓存里有数据
func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	// 数据不存在，err = redis.nil
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
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
