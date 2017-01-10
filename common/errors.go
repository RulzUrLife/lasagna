package common

import (
	"net/http"
)

type HTTPError struct {
	Code       int    `json:"status_code"`
	StatusText string `json:"message"`
	Msg        string `json:"payload"`
}

func NewHTTPError(msg string, code int) *HTTPError {
	return &HTTPError{code, http.StatusText(code), msg}
}

func New404Error(msg string) *HTTPError {
	return NewHTTPError(msg, http.StatusNotFound)
}

func New500Error(msg string) *HTTPError {
	return NewHTTPError(msg, http.StatusInternalServerError)
}
