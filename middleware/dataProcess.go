package middleware

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"go-do/common/dto"
	"strconv"
	"strings"
)

func DataProcess() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			resp, err := next(ctx, request)

			if err != nil {
				return nil, getErrorCode(err.Error())
			}
			return dto.NewSuccessResponseDto(resp), nil
		}
	}
}

func getErrorCode(errorStr string) *dto.ResponseDto {
	arr := strings.Split(errorStr, ",")
	intCode, err := strconv.ParseInt(arr[1], 10, 32)
	if err != nil {
		panic(err)
	}
	return dto.NewResponseDto(false, dto.DtoStatusCode(int(intCode)), arr[0], nil)
}
