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
	ErrUserDuplicate = errors.New("邮箱冲突 or 手机号码冲突")
	ErrUserNotFount  = gorm.ErrRecordNotFound
)

type UserDAO interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	Insert(ctx context.Context, u User) error
	FindByWechat(ctx context.Context, openID string) (User, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	//res := &GORMUserDAO{
	//		db: db,
	//	}
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	db := gorm.Open(mysql.Open())
	//	pt := unsafe.Pointer(&res.db)
	//	atomic.StorePointer(&pt, unsafe.Pointer(&db))
	//})
	//return res
	return &GORMUserDAO{
		db: db,
	}
}

func (dao *GORMUserDAO) FindByWechat(ctx context.Context, openID string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("wechat_open_id = ?", openID).First(&u).Error
	//方法二：err := dao.db.WithContext(ctx).First(&u,"email = ?", email).Error
	return u, err
}

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	//方法二：err := dao.db.WithContext(ctx).First(&u,"email = ?", email).Error
	return u, err
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("`id` = ?", id).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	// 存毫秒
	now := time.Now().UnixMilli()
	// SELECT * FROM users where email = 123@qq.com FOR UPDATE // 锁住间隙
	u.Utime = now
	u.Ctime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
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

	// 唯一索引允许有多个空值，即NULL
	// 但是不能有多个""
	Phone sql.NullString `gorm:"unique"`
	// 最大问题就是，你要解引用
	// 你要判空
	//Phone *string
	// 往这里面加

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
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}
