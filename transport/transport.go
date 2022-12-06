package transport

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-do/endpoint"
)

func ToTransport(handle func(ctx context.Context, params interface{}) (interface{}, error), params interface{}) gin.HandlerFunc {
	return NewServer(endpoint.ToEndpoint(handle), params)
}
