package custom

import (
	"reflect"
)

type ApiResponse[T any] struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    T      `json:"data"`
	Length  int    `json:"length"`
}

type ApiResponseWithPaginate[T any, P any] struct {
	Error        bool   `json:"error"`
	Message      string `json:"message"`
	Data         T      `json:"data"`
	PaginateData P      `json:"paginate"`
	Length       int    `json:"length"`
}

func Null() interface{} {
	return nil
}

func BuildResponse_[T any](status bool, message string, data T) ApiResponse[T] {
	var length int
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Array || val.Kind() == reflect.Slice {
		length = val.Len()
	} else {
		length = 0
	}

	return ApiResponse[T]{
		Error:   status,
		Message: message,
		Data:    data,
		Length:  length,
	}
}

func BuildResponse[T any](responseStatus ResponseStatus, data T) ApiResponse[T] {
	return BuildResponse_(responseStatus.GetResponseStatus(), responseStatus.GetResponseMessage(), data)
}

func BuildResponseWithPaginate_[T any, P any](status bool, message string, data T, paginateData P) ApiResponseWithPaginate[T, P] {
	var length int
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Array || val.Kind() == reflect.Slice {
		length = val.Len()
	} else {
		length = 0
	}

	return ApiResponseWithPaginate[T, P]{
		Error:        status,
		Message:      message,
		Data:         data,
		PaginateData: paginateData,
		Length:       length,
	}
}

func BuildResponseWithPaginate[T any, P any](responseStatus ResponseStatus, data T, paginateData P) ApiResponseWithPaginate[T, P] {
	return BuildResponseWithPaginate_(responseStatus.GetResponseStatus(), responseStatus.GetResponseMessage(), data, paginateData)
}
