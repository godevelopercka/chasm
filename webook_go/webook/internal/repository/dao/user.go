package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate = errors.New("邮箱或手机号码冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Save(ctx context.Context, id int64, Nickname, Birthday, AboutMe string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	Profile(ctx context.Context, id int64) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	Insert(ctx context.Context, u User) error
	FindByWechat(ctx context.Context, openID string) (User, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

// 编辑功能 Edit 插入新增字段
func (dao *GORMUserDAO) Save(ctx context.Context, id int64, Nickname, Birthday, AboutMe string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("Id = ?", id).First(&u).Error
	if err == nil {
		u.Nickname = Nickname
		u.Birthday = Birthday
		u.AboutMe = AboutMe
		dao.db.WithContext(ctx).Save(&u)
		return u, err
	}
	return u, err
}

func (dao *GORMUserDAO) FindByWechat(ctx context.Context, openID string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("wechat_open_id = ?", openID).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("`id` = ?", id).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) Profile(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("Id = ?", id).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	// 存毫秒
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062 // 数据库唯一索引冲突
		if mysqlErr.Number == 1062 {
			// 邮箱冲突 or 手机号码冲突
			return ErrUserDuplicate
		}
	}
	return err
}

type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	Nickname string
	Birthday string
	AboutMe  string
	Phone    sql.NullString `gorm:"unique"`
	// 最大问题是，你要解引用，接引用就要判空
	//Phone *string

	// 索引的最左匹配原则：
	// 假如索引在 <A, B, C>建好了
	// A, AB, ABC 都能用
	// WHERR A = ?
	// WHERR A = ? AND B = ? WHERR B = ? AND A = ?
	// WHERR A = ? AND B = ? AND C = ? 顺序随便换
	// WHERE 里面带了 ABC，可以用
	// WHERE 里面，没有 A，就不能用
	// 如果要创建联合索引，<unionid, openid>
	// <openid, unionid> 用 unionid 查询的时候，不会走索引
	// 微信的字段
	WechatUnionID sql.NullString
	WechatOpenID  sql.NullString `gorm:"unique"`

	Ctime int64
	Utime int64
}
