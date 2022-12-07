package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	"go-do/middleware"
	"reflect"
)

var dataProcess = middleware.DataProcess()

func Scan(d interface{}, ctx *gin.Context) any {
	v := reflect.ValueOf(d).Elem()
	iv := reflect.Indirect(v)
	s := iv.Type()
	for i := 0; i < s.NumField(); i++ {
		queryName := s.Field(i).Tag.Get("query")
		val := ctx.Query(queryName)
		v.Field(i).SetString(val)
	}
	return d
}

func query(t reflect.Type, ctx *gin.Context) any {
	value := reflect.New(t).Interface()
	return Scan(value, ctx)
}

func getMethodQuery(req interface{}, ctx *gin.Context) any {
	return query(reflect.TypeOf(req).Elem(), ctx)
}

// 1. decode request      http.request -> rpcModel.request
func decodeRequest(ctx *gin.Context, req interface{}) (interface{}, error) {
	if req == nil {
		return req, nil
	}
	if ctx.Request.Method == "GET" {
		req = getMethodQuery(req, ctx)
		return req, nil
	}
	err := ctx.ShouldBind(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// 2. encode response      rpcModel.response -> http.response
func encodeJsonResponse(ctx *gin.Context, res interface{}) error {
	ctx.JSON(200, res)
	return nil
}

//func decodeQueryRequest[T interface{}](ctx *gin.Context, req *T) (*T, error) {
//	err := ctx.Query(req)
//	if err != nil {
//		return nil, err
//	}
//	return req, nil
//}

func NewServer(endpoint endpoint.Endpoint, req interface{}) gin.HandlerFunc {
	hystrix := middleware.Hystrix()
	endpoint = hystrix(dataProcess(endpoint))
	return func(ctx *gin.Context) {
		req, err := decodeRequest(ctx, req)
		if err != nil {
			ctx.JSON(500, err)
			return
		}
		res, err := endpoint(ctx, req)
		if err != nil {
			ctx.JSON(500, err)
			return
		}
		err = encodeJsonResponse(ctx, res)
		if err != nil {
			ctx.JSON(500, err)
			return
		}
	}
}
