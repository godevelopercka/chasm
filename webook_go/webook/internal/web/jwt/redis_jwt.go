package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

var (
	AtKey = []byte("NDIOaqI8vCUZfWoNVcol0CuqFwHbu4cn")
	RtKey = []byte("NDIOaqI8vCUZfWoNVcol0CuqFwHbu4cc")
)

type RedisJWTHandler struct {
	// access_token key
	atKey []byte
	// refresh_token key
	rtKey []byte
	cmd   redis.Cmdable
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd: cmd,
	}
}

func (h *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = h.SetRefreshToken(ctx, uid, ssid)
	return err
}

func (h *RedisJWTHandler) SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := RefreshClaims{
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 过期时间：获取当前时间再加上1分钟
		},
		Uid: uid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(RtKey) // token 加密
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr) // 放到前端的 header 里面
	return nil
}

func (h *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	claims := ctx.MustGet("claims").(*UserClaims)
	return h.cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", claims.Ssid), "", time.Hour*24*7).Err()
}

func (h *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	_, err := h.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	return err
}

func (h *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	// 我现在用 JWT 来校验
	tokenHeader := ctx.GetHeader("Authorization")
	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

func (h RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)), // 过期时间：获取当前时间再加上1分钟
		},
		Id:        uid,
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(), // 拿到浏览器的 UserAgent, 也可以记录前端当时登录的设备信息，浏览信息等，然后打包传进来，这样可以保护 JWT 被盗用
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(AtKey) // token 加密
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr) // 放到前端的 header 里面
	return nil
}
