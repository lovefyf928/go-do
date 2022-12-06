package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-do/common/conf"
	"go-do/nacos"
	"go-do/register"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"net/rpc"
)

func Start(configPath string, serverPrefix string, registrar func(r *register.Register)) {

	errc := make(chan error)

	err := conf.LoadConfigInformation(configPath)

	nacos.LoadNacos()

	if err != nil {
		panic(err)
	}

	grpcArr := ":" + conf.ConfigInfo.Server.GrpcPort
	rpcArr := ":" + conf.ConfigInfo.Server.RpcPort

	engine := gin.New()

	grpcServer := grpc.NewServer()

	fmt.Println("grpc server running at: " + grpcArr)

	rpcServer := rpc.NewServer()

	r := register.NewRegister(engine, serverPrefix, grpcServer, rpcServer)

	registrar(r)

	r.Done()

	go runGrpc(grpcServer, grpcArr, errc)
	go runRpc(rpcServer, rpcArr, errc)
	go runHttp(engine, errc)

	logrus.WithField("error", <-errc).Info("Exit")

}

func runGrpc(grpcServer *grpc.Server, grpcArr string, errc chan error) {
	lis, err := net.Listen("tcp", grpcArr)
	if err != nil {
		errc <- err
	}
	err = grpcServer.Serve(lis)
	if err != nil {
		errc <- err
	}
}

func runHttp(engine *gin.Engine, errc chan error) {
	err := engine.Run(":" + conf.ConfigInfo.Server.Port)
	if err != nil {
		errc <- err
	}
}

func runRpc(rpcServer *rpc.Server, rpcArr string, errc chan error) {
	rpcServer.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	fmt.Println("rpc server running at: " + rpcArr)
	err := http.ListenAndServe(rpcArr, nil)
	if err != nil {
		errc <- err
	}
}
