package web

import (
	"errors"
	"geek-hw-week4/internal/domain"
	"geek-hw-week4/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type UserHandler struct {
	emailRegExp    *regexp.Regexp
	passwordRegExp *regexp.Regexp
	userSvc        service.UserService
	codeSvc        service.CodeService
}

func (handler *UserHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/users")
	g.POST("/singup", handler.Signup)
	g.POST("/login", handler.LoginJWT)
	g.POST("/edit", handler.Edit)
	g.GET("/profile", handler.Profile)

	// 手机验证码登录相关功能
	g.POST("/login_sms/code/send", handler.SendSMSLoginCode)
	g.POST("/login_sms", handler.LoginSMS)
}

func (handler *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := handler.codeSvc.Verify(ctx, "login", req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码不对，请重新输入",
		})
		return
	}
	u, err := handler.userSvc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	handler.setJWTToken(ctx, u.Id)
	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})
}

func (handler *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 你这边可以校验 Req
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "请输入手机号码",
		})
		return
	}
	err := handler.codeSvc.Send(ctx, "login", req.Phone)
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case errors.Is(err, service.ErrCodeSendTooMany):
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "短信发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 补日志的
	}
}

func (handler *UserHandler) Signup(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := handler.emailRegExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不对")
		return
	}

	isPassword, err := handler.passwordRegExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含字母、数字、特殊字符，并且不少于八位")
		return
	}

	err = handler.userSvc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch {
	case err == nil:
		ctx.String(http.StatusOK, "注册成功")
	case errors.Is(err, service.ErrDuplicateEmail):
		ctx.String(http.StatusOK, "邮箱冲突，请换一个")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}

func (handler *UserHandler) LoginJWT(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := handler.userSvc.Login(ctx, req.Email, req.Password)
	switch {
	case err == nil:
		handler.setJWTToken(ctx, u.Id)
		ctx.String(http.StatusOK, "登录成功")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (handler *UserHandler) setJWTToken(ctx *gin.Context, uid int64) {
	uc := UserClaims{
		Uid:       uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			// 1 分钟过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	ctx.Header("x-jwt-token", tokenStr)
}

func (handler *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Nickname string `json:"nickname"`
		// YYYY-MM-DD
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	uc, ok := ctx.MustGet("user").(UserClaims)
	if !ok {
		//ctx.String(http.StatusOK, "系统错误")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 用户输入不对
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		//ctx.String(http.StatusOK, "系统错误")
		ctx.String(http.StatusOK, "生日格式不对")
		return
	}
	err = handler.userSvc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:       uc.Uid,
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "更新成功")
}

func (handler *UserHandler) Profile(ctx *gin.Context) {
	uc, ok := ctx.MustGet("user").(UserClaims)
	if !ok {
		//ctx.String(http.StatusOK, "系统错误")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	u, err := handler.userSvc.FindById(ctx, uc.Uid)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	type User struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		AboutMe  string `json:"aboutMe"`
		Birthday string `json:"birthday"`
	}
	ctx.JSON(http.StatusOK, User{
		Nickname: u.Nickname,
		Email:    u.Email,
		AboutMe:  u.AboutMe,
		Birthday: u.Birthday.Format(time.DateOnly),
	})
}

func NewUserHandler(userSvc service.UserService, codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		emailRegExp:    regexp.MustCompile("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$", regexp.None),
		passwordRegExp: regexp.MustCompile(`^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`, regexp.None),
		userSvc:        userSvc,
		codeSvc:        codeSvc,
	}
}
