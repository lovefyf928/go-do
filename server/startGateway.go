package server

import (
	"github.com/gin-gonic/gin"
	"go-do/common/authorization"
	"go-do/common/conf"
	"go-do/middleware"
	"go-do/middleware/chain"
	"go-do/nacos"
)

func StartGateway(configPath string, fn func(chain *chain.Chain)) {
	err := conf.LoadConfigInformation(configPath)
	if err != nil {
		panic(err)
	}
	nacos.LoadNacos()
	authorization.LoadJwtConfig()

	if len(conf.ConfigInfo.GateWay.Routers) == 0 {
		panic("网关未配置路由，请检查")
	}

	if conf.ConfigInfo.Server.GatewayPort == "" {
		panic("请配置网关端口")
	}

	Run(":"+conf.ConfigInfo.Server.GatewayPort, fn)

}

func setupGin(fn func(chain *chain.Chain)) *gin.Engine {
	engine := gin.New()
	mw := chain.ChainMiddleware(fn)
	engine.Use(mw)
	engine.Use(middleware.Proxy)
	return engine

}

func Run(addr string, fn func(chain *chain.Chain)) {
	engine := setupGin(fn)
	err := engine.Run(addr)
	if err != nil {
		panic(err)
	}
}
