package domain

import (
	"time"
)

type User struct {
	Id       int64
	Email    string
	Password string
	Nickname string
	Birthday string
	AboutMe  string
	Phone    string
	// 不要组合，万一你可能还有钉钉的相同字段 UionID
	WechatInfo WechatInfo
	Ctime      time.Time
}
