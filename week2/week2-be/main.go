package main

import (
	"fzy.com/geek-hw-week2/controller"
	"fzy.com/geek-hw-week2/controller/middleware"
	"fzy.com/geek-hw-week2/repository"
	"fzy.com/geek-hw-week2/repository/dao"
	"fzy.com/geek-hw-week2/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func main() {
	db := initDB()
	server := initServer()
	initWebControllers(db, server)
	server.Run(":8080")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:admin111@tcp(localhost:3306)/geek_hw"))
	if err != nil {
		panic("fail to connect database")
	}

	if err := dao.InitTables(db); err != nil {
		panic("fail to init tables")
	}

	return db
}

func initServer() *gin.Engine {
	server := gin.Default()

	// configure CORS
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	// configure session
	server.Use(sessions.Sessions("sid", cookie.NewStore([]byte("geek_hw"))))

	// configure login middleware
	loginMiddleware := &middleware.LoginMiddlewareBuilder{}
	server.Use(loginMiddleware.CheckLogin())

	return server
}

func initWebControllers(db *gorm.DB, server *gin.Engine) {
	accountDAO := dao.NewAccountDAO(db)
	accountRepository := repository.NewAccountRepository(accountDAO)
	accountService := service.NewAccountService(accountRepository)
	accountController := controller.NewAccountController(accountService)
	accountController.RegisterRoutes(server)
}