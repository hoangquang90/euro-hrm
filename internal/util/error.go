package util

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
)

// NewError example
func NewError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Code: status,
	}
	var ve validator.ValidationErrors
	var pgErr *pgconn.PgError
	if errors.As(err, &ve) {
		details := make([]DetailError, len(ve))
		for i, fe := range ve {
			message := ValidationErrorToText(&fe)
			details[i] = DetailError{fe.Field(), fe.Value(), message}
			er.Message += message + "\n"
		}
		er.Details = details
		er.Message = strings.Trim(er.Message, "\n")
	} else if errors.As(err, &pgErr) {
		er.Message = pgErr.Message
	} else {
		er.Message = err.Error()
	}
	ctx.AbortWithStatusJSON(status, er)
}

// HTTPError example
type HTTPError struct {
	Code    int           `json:"code" example:"400"`
	Message string        `json:"message" example:"status bad request"`
	Details []DetailError `json:"details"`
}

type DetailError struct {
	Param   string      `json:"param"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

func ValidationErrorToText(e *validator.FieldError) string {
	field := (*e).Field()
	param := (*e).Param()
	kind := (*e).Kind()
	switch (*e).Tag() {
	case "len=0|alphanum":
		return fmt.Sprintf("%s must be alphabet or number or not required", field)
	case "required", "notblank":
		return fmt.Sprintf("%s is required", field)
	case "max":
		if kind == reflect.Int {
			return fmt.Sprintf("%s maximum value is %s", field, param)
		} else {
			return fmt.Sprintf("%s maximum length is %s", field, param)
		}
	case "min":
		if kind == reflect.Int {
			return fmt.Sprintf("%s minimum value is %s", field, param)
		} else {
			return fmt.Sprintf("%s minimum length is %s", field, param)
		}
	case "email":
		return "Invalid email format"
	case "len":
		return fmt.Sprintf("%s must be %s characters long", field, param)
	case "url":
		return "Invalid URL format"
	}
	return fmt.Sprintf("%s is not valid", field)
}

// func TranslateValidatorError(err error, code int) HTTPError {
// 	result := HTTPError{
// 		Code: code,
// 	}
// 	var detail []DetailError
// 	var ve validator.ValidationErrors
// 	if errors.As(err, &ve) {
// 		detail = make([]DetailError, len(ve))
// 		for i, fe := range ve {
// 			message := ValidationErrorToText(&fe)
// 			detail[i] = DetailError{fe.Field(), fe.Value(), message}
// 			result.Message += message + "\n"
// 		}
// 	}
// 	result.Details = detail
// 	return result
// }
