package middleware

import (
	"context"
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
)

func Hystrix() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			ctx1, ok := ctx.(*gin.Context)
			if ok {
				var resp interface{}
				if err := hystrix.Do(ctx1.Request.URL.Path, func() (err error) {
					resp, err = next(ctx, request)
					return err
				}, nil); err != nil {
					return nil, err
				}
				return resp, nil
			}
			return nil, errors.New("context format error")
		}
	}
}
