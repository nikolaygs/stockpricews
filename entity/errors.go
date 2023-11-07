package entity

import "errors"

var ErrBadRequest = errors.New("bad request")
var ErrNotFound = errors.New("not found")
var ErrMethodNotAllowed = errors.New("method not allowed")

type ErrorMessage struct {
	Message string `json:"message"`
}
