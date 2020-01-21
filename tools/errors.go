package tools

import "errors"

var (
	ErrAlreadyExists  = errors.New("already exists")
	ErrNotFound       = errors.New("not found")
	ErrHTTPBadRequest = errors.New("bad request")
)
