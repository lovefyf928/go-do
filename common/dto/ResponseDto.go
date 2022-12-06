package dto

import "encoding/json"

type ResponseDto struct {
	Success    bool          `json:"success"`
	StatusCode DtoStatusCode `json:"statusCode"`
	Msg        string        `json:"msg"`
	Data       interface{}   `json:"data"`
}

func NewResponseDto(success bool, statusCode DtoStatusCode, msg string, data interface{}) *ResponseDto {
	return &ResponseDto{Success: success, StatusCode: statusCode, Msg: msg, Data: data}
}

func NewSuccessResponseDto(data interface{}) *ResponseDto {
	return &ResponseDto{Success: true, StatusCode: SUCCESS, Msg: "成功", Data: data}
}

func NewSuccessResponseDtoNilMsg(msg string) *ResponseDto {
	return &ResponseDto{Success: true, StatusCode: SUCCESS, Msg: msg, Data: nil}
}

func (r ResponseDto) Error() string {
	marshal, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(marshal)
}

type DtoStatusCode int

const (
	SUCCESS         DtoStatusCode = 200
	INTERNATL_ERROR DtoStatusCode = 500
	FORBBDIEN       DtoStatusCode = 401
	UNAUTHORIZED    DtoStatusCode = 403
)
