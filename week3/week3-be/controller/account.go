package controller

import (
	"errors"
	"net/http"
	"time"

	"fzy.com/geek-hw-week3/controller/middleware"
	"fzy.com/geek-hw-week3/domain"
	"fzy.com/geek-hw-week3/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	emailRegexPattern    = "^((?!\\.)[\\w-_.]*[^.])(@\\w+)(\\.\\w+(\\.\\w+)?[^.\\W])$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type AccountController struct {
	emailRegExp    *regexp.Regexp
	passwordRegExp *regexp.Regexp
	accountService *service.AccountService
}

func NewAccountController(accountService *service.AccountService) *AccountController {
	return &AccountController{
		emailRegExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		accountService: accountService,
	}
}

func (accountController *AccountController) RegisterRoutes(server *gin.Engine) {
	routes := server.Group("/account")
	routes.POST("/signup", accountController.Signup)
	//routes.POST("/login", accountController.LoginWithSession)
	routes.POST("/login", accountController.LoginWithJWT)
	routes.PUT("/edit", accountController.Edit)
	routes.GET("/profile", accountController.Profile)
}

func (accountController *AccountController) Signup(ctx *gin.Context) {
	type SignupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignupReq
	if err := ctx.Bind(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	// validate email
	isEmail, err := accountController.emailRegExp.MatchString(req.Email)
	if err != nil || !isEmail {
		ctx.Status(http.StatusBadRequest)
		return
	}

	// validate password
	if req.Password != req.ConfirmPassword {
		ctx.Status(http.StatusBadRequest)
		return
	}
	isPassword, err := accountController.passwordRegExp.MatchString(req.Password)
	if err != nil || !isPassword {
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = accountController.accountService.Signup(ctx, domain.Account{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrDuplicateEmail) {
		ctx.Status(http.StatusBadRequest)
		return
	}
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func (accountController *AccountController) LoginWithSession(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	account, err := accountController.accountService.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidEmailOrPassword) {
		ctx.Status(http.StatusBadRequest)
		return
	}
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	s := sessions.Default(ctx)
	s.Set("account_id", account.Id)
	s.Options(sessions.Options{
		MaxAge: 900,
	})
	if err := s.Save(); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func (accountController *AccountController) LoginWithJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	account, err := accountController.accountService.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidEmailOrPassword) {
		ctx.Status(http.StatusBadRequest)
		return
	}
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, middleware.AccountJWTClaims{
		AccountId: account.Id,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	})
	tokenStr, err := token.SignedString(middleware.JWTSecret)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Header("x-jwt", tokenStr)
	ctx.Status(http.StatusOK)
}

func (accountController *AccountController) Edit(ctx *gin.Context) {
	type EditReq struct {
		Name  string `json:"name"`
		Birth string `json:"birth"`
		About string `json:"about"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	//s := sessions.Default(ctx)
	//if s.Get("account_id") == nil {
	//	ctx.Status(http.StatusNotFound)
	//	return
	//}
	//accountId := s.Get("account_id").(int64)

	claims, ok := ctx.Get("account")
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	var accountClaims middleware.AccountJWTClaims
	accountClaims = claims.(middleware.AccountJWTClaims)
	account, err := accountController.accountService.GetProfileById(ctx, accountClaims.AccountId)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	if err := accountController.accountService.EditProfile(ctx, account, req.Name, req.Birth, req.About); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func (accountController *AccountController) Profile(ctx *gin.Context) {
	//s := sessions.Default(ctx)
	//if s.Get("account_id") == nil {
	//	ctx.Status(http.StatusNotFound)
	//	return
	//}
	//accountId := s.Get("account_id").(int64)

	claims, ok := ctx.Get("account")
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	var accountClaims middleware.AccountJWTClaims
	accountClaims = claims.(middleware.AccountJWTClaims)
	account, err := accountController.accountService.GetProfileById(ctx, accountClaims.AccountId)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, account)
}
