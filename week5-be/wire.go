//go:build wireinject

package main

import (
	"geek-hw-week4/internal/repository"
	"geek-hw-week4/internal/repository/cache"
	"geek-hw-week4/internal/repository/dao"
	"geek-hw-week4/internal/service"
	"geek-hw-week4/internal/web"
	"geek-hw-week4/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,

		dao.NewUserDAO,

		cache.NewUserCache, cache.NewCodeCache,

		repository.NewUserRepository, repository.NewCodeRepository,

		ioc.InitLocalSMSService,
		service.NewUserService, service.NewCodeService,

		web.NewUserHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
