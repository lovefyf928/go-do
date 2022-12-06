package routers

import (
	"github.com/gin-gonic/gin"
	"go-do/common/enum"
)

type IRouter interface {
	NewRouter(engine *gin.Engine, prefix string) IRouter

	AddRouter(method enum.HttpMethod, path string, fn func(*gin.Context))

	AddGroup(groupPrefix string)

	Register()
}

type RouterFactory struct {
	IRouter
}
