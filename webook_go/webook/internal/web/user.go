package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"webook_go/webook/internal/domain"
	"webook_go/webook/internal/service"
)

type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	BirthdayExp *regexp.Regexp
	svc         *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	// 定义校验邮箱和密码的正则表达式
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,72}$`
		birthRegexPattern    = `\b\d{4}-\d{2}-\d{2}\b`
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
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
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
	// 设置好 session 必须要保存才能生效
	sess.Save()

	ctx.String(http.StatusOK, "登录成功")
	fmt.Println(user)
	return
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
