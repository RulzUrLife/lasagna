package common

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	Code       int    `json:"status_code"`
	StatusText string `json:"message"`
	Msg        string `json:"payload"`
}

func NewHTTPError(code int, msg string, i ...interface{}) *HTTPError {
	return &HTTPError{code, http.StatusText(code), fmt.Sprintf(msg, i...)}
}

func New404Error(msg string, i ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusNotFound, msg, i...)
}

func New500Error(msg string, i ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusInternalServerError, msg, i...)
}

func New400Error(msg string, i ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusBadRequest, msg, i...)
}
