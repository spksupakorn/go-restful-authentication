package custom

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func PanicController(c *gin.Context) {
	if err := recover(); err != nil {
		str := fmt.Sprint(err)
		strArr := strings.Split(str, ":")

		key := strArr[0]
		msg := strings.Trim(strArr[1], " ")

		code, _ := strconv.Atoi(key)

		switch code {
		case
			http.StatusNotFound:
			c.JSON(http.StatusNotFound, BuildResponse_(DataNotFound.GetResponseStatus(), msg, Null()))
			c.Abort()
		case
			http.StatusBadRequest:
			c.JSON(http.StatusBadRequest, BuildResponse_(BadRequest.GetResponseStatus(), msg, Null()))
			c.Abort()
		case
			http.StatusUnauthorized:
			c.JSON(http.StatusUnauthorized, BuildResponse_(Unauthorized.GetResponseStatus(), msg, Null()))
			c.Abort()
		case
			http.StatusForbidden:
			c.JSON(http.StatusForbidden, BuildResponse_(Forbidden.GetResponseStatus(), msg, Null()))
			c.Abort()
		case
			http.StatusUnprocessableEntity:
			c.JSON(http.StatusUnprocessableEntity, BuildResponse_(UnprocessableEntity.GetResponseStatus(), msg, Null()))
			c.Abort()
		default:
			c.JSON(http.StatusInternalServerError, BuildResponse_(InternalServerError.GetResponseStatus(), msg, Null()))
			c.Abort()
		}
	}
}

func PanicException(err error) {
	if appErr, ok := err.(*AppError); ok {
		PanicException_(appErr.Code, appErr.Message)
	} else if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrs {
			feildErr := fmt.Sprintf("'%s' field is required", err.Field())
			PanicException_(http.StatusBadRequest, feildErr)
		}
		return
	} else {
		PanicException_(http.StatusInternalServerError, "unexpected error")
	}
}

func PanicException_(statusCode int, message string) {
	err := errors.New(message)
	err = fmt.Errorf("%d: %w", statusCode, err)
	if err != nil {
		// log.Error(err)
		panic(err)
	}
}
