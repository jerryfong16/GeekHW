package main

import (
	"fzy.com/geek-hw-week2/config"
	"fzy.com/geek-hw-week2/controller"
	"fzy.com/geek-hw-week2/controller/middleware"
	"fzy.com/geek-hw-week2/repository"
	"fzy.com/geek-hw-week2/repository/dao"
	"fzy.com/geek-hw-week2/service"
	"github.com/gin-contrib/cors"
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
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
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
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		ExposeHeaders: []string{"x-jwt"},
		AllowMethods:  []string{"OPTIONS", "POST", "DELETE", "PUT", "GET"},
		MaxAge:        12 * time.Hour,
	}))

	// configure login session middleware
	//server.Use(sessions.Sessions("sid", cookie.NewStore([]byte("geek_hw"))))
	//loginMiddleware := &middleware.LoginMiddlewareBuilder{}
	//server.Use(loginMiddleware.CheckLogin())

	// configure login jwt middleware
	loginJWTMiddleware := &middleware.LoginJWTMiddlewareBuilder{}
	server.Use(loginJWTMiddleware.CheckLogin())

	return server
}

func initWebControllers(db *gorm.DB, server *gin.Engine) {
	accountDAO := dao.NewAccountDAO(db)
	accountRepository := repository.NewAccountRepository(accountDAO)
	accountService := service.NewAccountService(accountRepository)
	accountController := controller.NewAccountController(accountService)
	accountController.RegisterRoutes(server)
}
