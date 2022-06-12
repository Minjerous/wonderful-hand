package errdef

import (
	"fmt"
)

type Err struct {
	HttpCode    int
	StatusCode  int
	Description string
}

var Nil = New(-1, -1, "")

func (e Err) Error() string {
	return fmt.Sprintf("HttpCode: %d StatusCode: %d Description: %s\n", e.HttpCode, e.StatusCode, e.Description)
}

func New(httpCode, statusCode int, des string) Err {
	return Err{
		HttpCode:    httpCode,
		StatusCode:  statusCode,
		Description: des,
	}
}

func Errorf(httpCode, statusCode int, format string, a ...any) Err {
	return Err{
		HttpCode:    httpCode,
		StatusCode:  statusCode,
		Description: fmt.Sprintf(format, a...),
	}
}

func IsNil(err Err) bool {
	return err == Nil || err.HttpCode < 100
}
