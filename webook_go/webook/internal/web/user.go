package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"webook_go/webook/internal/domain"
	"webook_go/webook/internal/service"
	ijwt "webook_go/webook/internal/web/jwt"
)

const biz = "login"

// 确保 UserHandler 上实现了 handler 接口
var _ handler = &UserHandler{}

// 这个更优雅
var _ handler = (*UserHandler)(nil)

type UserHandler struct {
	emailExp    *regexp.Regexp
	codeSvc     service.CodeService
	passwordExp *regexp.Regexp
	BirthdayExp *regexp.Regexp
	svc         service.UserService
	ijwt.Handler
	cmd redis.Cmdable
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService, jwtHdl ijwt.Handler) *UserHandler {
	// 定义校验邮箱和密码的正则表达式
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,72}$`
		//birthRegexPattern    = `\b\d{4}-\d{2}-\d{2}\b`
		birthRegexPattern = `^\d{4}-\d{2}-\d{2}$`
	)
	// 预编译正则表达式，提高校验速度
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	birthdayExp := regexp.MustCompile(birthRegexPattern, regexp.None)
	return &UserHandler{
		emailExp:    emailExp,
		passwordExp: passwordExp,
		BirthdayExp: birthdayExp,
		svc:         svc,
		codeSvc:     codeSvc,
		Handler:     jwtHdl,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	//ug.POST("/edit", u.Edit)
	ug.POST("/edit", u.EditJWT)
	//ug.GET("/profile", u.Profile)
	ug.GET("/profile", u.ProfileJWT)
	ug.GET("/logout", u.LogOutJWT)
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
	ug.POST("/refresh_token", u.RefreshToken)
}

// RefreshToken 可以同时刷新长短 token，用 redis 来记录是否有效，即 refresh_token 是一次性的
// 参考登录校验部分，比较 User-Agent 来增强安全性
func (u *UserHandler) RefreshToken(ctx *gin.Context) {
	// 只有这个接口，拿出来的才是 refresh_token，其他地方都是 access token
	refreshToken := u.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(refreshToken, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RtKey, nil
	})
	if err != nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = u.CheckSession(ctx, rc.Ssid)
	// 搞个新的 access_token
	err = u.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		// 这种系统异常的含糊日志不要打，信息量不足，无法明确知道发生了什么错误，要坐到能根据日志定位问题
		zap.L().Error("系统异常", zap.Error(err))
		// 正常来说，msg 的部分就应该包含足够的定位信息
		zap.L().Error("sdafagsaf 设置 JWT token 出现异常", zap.Error(err),
			zap.String("method", "UserHandler:RefreshToken"))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "刷新成功",
	})
}

func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 这边，可以加上各种校验
	ok, err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("校验验证码出错", zap.Error(err),
			// 不能这样打，因为手机号码是敏感数据，你不能打到日志里面
			// 打印加密后的
			zap.String("手机号码", req.Phone))
		// 最多最多就这样
		zap.L().Debug("", zap.String("手机号码", req.Phone))
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码有误",
		})
		return
	}

	// 我这个手机号会不会是一个新用户呢
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	if err = u.SetLoginToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 4,
		Msg:  "验证码校验通过",
	})
}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 是不是一个合法的手机号
	// 生产要换成正则表达式
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "输入有误",
		})
		return
	}
	err := u.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		zap.L().Error("发送太频繁", zap.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送太频繁，请稍候重试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq

	// bind 方法接收请求参数
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 校验邮箱
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式不正确")
		return
	}
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字、英文字母、特殊字符")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}
	// 调用一下 svc 的方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
	fmt.Printf("%v", req)
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if err = u.SetLoginToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "登录成功")
	fmt.Println(user)
	return
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 步骤二
	// 这里把 cookies 放到 session 中
	// 在这里登录成功后，要把 session 拿出来
	sess := sessions.Default(ctx)
	// 可以随便设置值了
	// 你要放在 session 里面的值
	// 设置好 session 后才能去校验
	sess.Set("userId", user.Id) // 这里把 Id 塞进了 session 中，所以登录校验只要证明存不存在这个 Id 就行
	sess.Options(sessions.Options{
		// 生产环境需要设置 Secure、 HttpOnly
		//Secure: true,
		//HttpOnly: true,
		MaxAge: 60, // 设置过期时间 30 秒
	})
	// 设置好 session 必须要保存才能生效
	sess.Save()

	ctx.String(http.StatusOK, "登录成功")
	fmt.Println(user)
	return
}

func (u *UserHandler) LogOutJWT(ctx *gin.Context) {
	err := u.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "退出登录失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "退出登录OK",
	})
}

func (u *UserHandler) LogOut(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		// 生产环境需要设置 Secure、 HttpOnly
		//Secure: true,
		//HttpOnly: true,
		MaxAge: -1, // 设置过期时间
	})
	// 设置好 session 必须要保存才能生效
	sess.Save()
	ctx.String(http.StatusOK, "退出登录成功")
}

func (u *UserHandler) EditJWT(ctx *gin.Context) {
	c, ok := ctx.Get("claims")
	// 你可以断定，必然有 claims
	//if !ok {
	//	// 你可以考虑监控住这里
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	// ok 代表是不是 *UserClaims
	claims, ok := c.(*ijwt.UserClaims) // 断言 c 是不是 UserClaims 的指针，如果不是就会 panic
	if !ok {
		// 你可以考虑监控住这里
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	type EditReq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 校验生日
	ok, err := u.BirthdayExp.MatchString(req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "生日日期格式不对")
		return
	}
	fmt.Println(claims.Uid)
	// 调用一下 svc 的方法
	user, err := u.svc.Edit(ctx, claims.Uid, req.Nickname, req.Birthday, req.AboutMe)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	if len(req.Nickname) > 24 {
		ctx.String(http.StatusOK, "昵称过长")
		return
	}
	if len(req.AboutMe) > 1024 {
		ctx.String(http.StatusOK, "个人简介过长")
		return
	}
	ctx.String(http.StatusOK, "提交成功")
	fmt.Printf("%v", user)
	fmt.Printf("%v", req)
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	// 获取已登录的 sessionId,不然不允许访问该路径
	sess := sessions.Default(ctx)
	id := sess.Get("userId").(int64)
	if id == 0 {
		// 没有登录
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	type EditReq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 校验生日
	ok, err := u.BirthdayExp.MatchString(req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "生日日期格式不对")
		return
	}
	// 调用一下 svc 的方法
	user, err := u.svc.Edit(ctx, id, req.Nickname, req.Birthday, req.AboutMe)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	if len(req.Nickname) > 24 {
		ctx.String(http.StatusOK, "昵称过长")
		return
	}
	if len(req.AboutMe) > 1024 {
		ctx.String(http.StatusOK, "个人简介过长")
		return
	}
	ctx.String(http.StatusOK, "提交成功")
	fmt.Printf("%v", user)
	fmt.Printf("%v", req)
}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, ok := ctx.Get("claims")
	// 你可以断定，必然有 claims
	//if !ok {
	//	// 你可以考虑监控住这里
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	// ok 代表是不是 *UserClaims
	claims, ok := c.(*ijwt.UserClaims) // 断言 c 是不是 UserClaims 的指针，如果不是就会 panic
	if !ok {
		// 你可以考虑监控住这里
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	type Profile struct {
		Email    string
		Nickname string
		Birthday string
		AboutMe  string
	}
	user, err := u.svc.Profile(ctx, claims.Uid)
	if err != nil {
		// 按照道理来说，这边 id 对应的数据肯定存在，所以要是没找到，
		// 那就说明是系统出了问题。
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Profile{
		Email:    user.Email,
		Nickname: user.Nickname,
		Birthday: user.Birthday,
		AboutMe:  user.AboutMe,
	})
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	id := sess.Get("userId").(int64)
	type Profile struct {
		Email    string
		Nickname string
		Birthday string
		AboutMe  string
	}
	user, err := u.svc.Profile(ctx, id)
	if err != nil {
		// 按照道理来说，这边 id 对应的数据肯定存在，所以要是没找到，
		// 那就说明是系统出了问题。
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Profile{
		Email:    user.Email,
		Nickname: user.Nickname,
		Birthday: user.Birthday,
		AboutMe:  user.AboutMe,
	})
}
