// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"geek-hw-week4/internal/repository"
	"geek-hw-week4/internal/repository/cache"
	"geek-hw-week4/internal/repository/dao"
	"geek-hw-week4/internal/service"
	"geek-hw-week4/internal/web"
	"geek-hw-week4/ioc"
	"github.com/gin-gonic/gin"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := ioc.InitRedis()
	v := ioc.InitGinMiddlewares(cmdable)
	db := ioc.InitDB()
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	userHandler := web.NewUserHandler(userService)
	engine := ioc.InitWebServer(v, userHandler)
	return engine
}
