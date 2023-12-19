package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	"webook_go/webook/internal/domain"
	"webook_go/webook/internal/service"
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
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
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
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
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

	// 这边要怎么办呢
	// 从哪来
	if err = u.setJWTToken(ctx, user.Id); err != nil {
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
	// 步骤二
	// 这里生成一个 JWT 设置登录态
	// 生成一个 JWT token

	if err = u.setJWTToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "登录成功")
	fmt.Println(user)
	return
}

func (u *UserHandler) setJWTToken(ctx *gin.Context, uid int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)), // 过期时间：获取当前时间再加上1分钟
		},
		Uid:       uid,
		UserAgent: ctx.Request.UserAgent(), // 拿到浏览器的 UserAgent, 也可以记录前端当时登录的设备信息，浏览信息等，然后打包传进来，这样可以保护 JWT 被盗用
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("NDIOaqI8vCUZfWoNVcol0CuqFwHbu4cn")) // token 加密
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr) // 放到前端的 header 里面
	return nil
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
	claims, ok := c.(*UserClaims) // 断言 c 是不是 UserClaims 的指针，如果不是就会 panic
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
	claims, ok := c.(*UserClaims) // 断言 c 是不是 UserClaims 的指针，如果不是就会 panic
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

type UserClaims struct {
	jwt.RegisteredClaims
	// 声明你自己的要放进去 token 里面的数据
	Uid       int64
	UserAgent string
}
