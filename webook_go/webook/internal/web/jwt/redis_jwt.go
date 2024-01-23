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

// AtKey 和 RtKey 是两个用于JWT加密的密钥，它们分别用于access_token和refresh_token
var (
	AtKey = []byte("NDIOaqI8vCUZfWoNVcol0CuqFwHbu4cn")
	RtKey = []byte("NDIOaqI8vCUZfWoNVcol0CuqFwHbu4cc")
)

type RedisJWTHandler struct {
	// access_token key
	atKey []byte
	// refresh_token key
	rtKey []byte
	// 用于与Redis交互的命令对象
	cmd redis.Cmdable
}

// NewRedisJWTHandler 用于创建一个新的RedisJWTHandler实例
func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd: cmd,
	}
}

// SetLoginToken 方法用于设置登录令牌。它首先生成一个唯一的ssid，然后使用这个ssid设置JWT Token，最后设置Refresh Token
func (h *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()          // 生成一个唯一的ssid
	err := h.SetJWTToken(ctx, uid, ssid) // 设置JWT Token
	if err != nil {
		return err
	}
	err = h.SetRefreshToken(ctx, uid, ssid)
	return err
}

// SetRefreshToken 方法用于设置Refresh Token。它首先生成一个带有uid和ssid的Refresh Claims
// 然后使用这些Claims生成一个新的JWT Token，最后将这个Token放入请求的Header中
func (h *RedisJWTHandler) SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := RefreshClaims{
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 过期时间：获取当前时间再加上1分钟
		},
		Uid: uid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) // 使用HS512签名方法创建一个新的JWT Token
	tokenStr, err := token.SignedString(RtKey)                 // token 加密
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr) // 放到前端的 header 里面
	return nil
}

// 清除用户Token的方法
// 从上下文的Header中删除JWT和Refresh Token
// 从Redis中删除用户的会话数据
func (h *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")                 // 清除JWT Token
	ctx.Header("x-refresh-token", "")             // 清除Refresh Token
	claims := ctx.MustGet("claims").(*UserClaims) // 从上下文中获取用户声明
	// 删除Redis中的会话数据
	return h.cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", claims.Ssid), "", time.Hour*24*7).Err()
}

// 检查用户会话的方法。
// 通过给定的会话ID检查Redis中是否存在该会话数据
func (h *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	_, err := h.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result() // 检查Redis中是否存在该会话数据
	return err                                                               // 返回错误（如果有）或nil（如果存在会话）
}

// 从上下文的Header中提取JWT Token。
// 使用Authorization头部的第一个值作为JWT Token
func (h *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	// 我现在用 JWT 来校验
	tokenHeader := ctx.GetHeader("Authorization")
	segs := strings.Split(tokenHeader, " ") // 根据空格分割Header值，期望的格式是 "Bearer [token]"
	if len(segs) != 2 {                     // 如果分割后的数组长度不为2，说明格式不正确，返回空字符串
		return ""
	}
	return segs[1] // 返回实际的JWT Token值
}

// 设置用户JWT Token的方法
// 生成一个新的JWT，并使用给定的密钥对其进行签名
// 将签名的JWT放入上下文的Header中
func (h RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{ // 创建用户声明，包括标准的和自定义的声明
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)), // 过期时间：获取当前时间再加上30分钟
		},
		Id:        uid,                     // 用户ID作为自定义声明的一部分，可以在JWT验证时使用此ID进行验证或授权决策
		Ssid:      ssid,                    // 会话ID作为自定义声明的一部分，用于识别特定会话
		UserAgent: ctx.Request.UserAgent(), // 拿到浏览器的 UserAgent, 也可以记录前端当时登录的设备信息，浏览信息等，然后打包传进来，这样可以保护 JWT 被盗用
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) // 使用HS512签名方法创建一个新的JWT，并使用给定的声明和密钥
	tokenStr, err := token.SignedString(AtKey)                 // token 加密
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr) // 放到前端的 header 里面
	return nil                          // 返回nil表示成功
}
