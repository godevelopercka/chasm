package repository

import (
	"context"
	"database/sql"
	"time"
	"webook_go/webook/internal/domain"
	"webook_go/webook/internal/repository/cache"
	"webook_go/webook/internal/repository/dao"
)

var ErrUserDuplicate = dao.ErrUserDuplicate
var ErrUserNotFound = dao.ErrUserNotFound
var ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes

type UserRepository interface {
	FindById(ctx context.Context, id int64) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	Edit(ctx context.Context, id int64, Nickname, Birthday, AboutMe string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindByWechat(ctx context.Context, openID string) (domain.User, error)
}

type CacheUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, c cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *CacheUserRepository) FindByWechat(ctx context.Context, openID string) (domain.User, error) {
	u, err := r.dao.FindByWechat(ctx, openID)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CacheUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 先从 cache 里面找
	// 再从 dao 里面找
	// 注册后 domain.user 的值就是数据库的 User 的，所以 redis 直接返回 domain.User 就是有值的
	// redis 就是缓存，有没有数据只有自己知道
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		// 必然有数据
		return u, nil
	}
	// 没这个数据
	if err == cache.ErrKeyNotExist {
		//去数据库里面加载
	}
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = r.entityToDomain(ue)

	//_ = r.cache.Set(ctx, u)
	//if err != nil {
	//	// 我这里怎么办？
	//	// 打日志，做监控
	//	//return domain.User{}, err
	//}
	go func() {
		_ = r.cache.Set(ctx, u)
	}()
	return u, nil
}

func (r *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntiy(u))
}

func (r *CacheUserRepository) Edit(ctx context.Context, id int64, Nickname, Birthday, AboutMe string) (domain.User, error) {
	u, err := r.dao.Save(ctx, id, Nickname, Birthday, AboutMe)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
	}, nil
}

func (r *CacheUserRepository) Profile(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.dao.Profile(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CacheUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
		Phone:    u.Phone.String,
		WechatInfo: domain.WechatInfo{
			UnionID: u.WechatUnionID.String,
			OpenID:  u.WechatOpenID.String,
		},
		Ctime: time.UnixMilli(u.Ctime),
	}
}

func (r *CacheUserRepository) domainToEntiy(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			// 我确实有手机号
			Valid: u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
		WechatOpenID: sql.NullString{
			String: u.WechatInfo.OpenID,
			Valid:  u.WechatInfo.OpenID != "",
		},
		WechatUnionID: sql.NullString{
			String: u.WechatInfo.UnionID,
			Valid:  u.WechatInfo.UnionID != "",
		},
		Ctime: u.Ctime.UnixMilli(),
	}
}
