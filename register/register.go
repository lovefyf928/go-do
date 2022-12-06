package register

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-do/common/enum"
	"go-do/common/utils/routers"
	"go-do/common/utils/routers/impl"
	"go-do/transport"
	"google.golang.org/grpc"
	"net/rpc"
)

type Register struct {
	router *routers.IRouter

	grpcServer *grpc.Server

	rpcServer *rpc.Server
}

func NewRegister(engine *gin.Engine, httpPrefix string, grpcServer *grpc.Server, rpcServer *rpc.Server) *Register {
	var router routers.IRouter = routers.RouterFactory{IRouter: &impl.GinRouter{}}

	router = router.NewRouter(engine, httpPrefix)

	return &Register{
		router:     &router,
		grpcServer: grpcServer,
		rpcServer:  rpcServer,
	}
}

func (r *Register) RegisterHttpHandle(path string, method enum.HttpMethod, handle func(ctx context.Context, req interface{}) (interface{}, error), params interface{}) {
	router := *r.router
	router.AddRouter(method, path, transport.ToTransport(handle, params))
}

func (r *Register) RegisterGrpcHandles(registerGrpcFn func(s grpc.ServiceRegistrar, svr any), implementedGrpcStruct any) {
	registerGrpcFn(r.grpcServer, implementedGrpcStruct)
}

func (r *Register) RegisterRpcHandles(name string, rcvr any) {
	err := r.rpcServer.RegisterName(name, rcvr)
	if err != nil {
		panic(err)
	}
}

func (r *Register) Done() {
	router := *r.router
	router.Register()
}
