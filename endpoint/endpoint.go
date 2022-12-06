package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func ToEndpoint(handle func(ctx context.Context, params interface{}) (interface{}, error)) endpoint.Endpoint {
	return handle
}
