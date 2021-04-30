package common

import (
	"iris_project_foundation/common/api_error"
)

type Response struct {
	Code    int64       `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
	Next    int64       `json:"next,omitempty"`
	Total   int64       `json:"total,omitempty"`
}

func Success(data interface{}, next, total int64) *Response {
	return &Response{
		Code:    10000,
		Message: "success",
		Data:    data,
		Next:    next,
		Total:   total,
	}
}

func Failed(e error) *Response {
	var (
		apiErr *api_error.BaseAPIError
		ok     bool
	)

	if apiErr, ok = e.(*api_error.BaseAPIError); ok {
		return &Response{Code: apiErr.ErrorCode, Message: apiErr.Error()}
	}
	return &Response{Code: 10001, Message: e.Error()}
}
