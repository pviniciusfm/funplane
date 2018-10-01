package server

import (
	"errors"
)

const (
	errInternalError = "internal error"
	errCache         = "cache error"
	errAdapter       = "adapter error"
	errNotFound      = "object not found"
)

var (
	ErrInternalError = errors.New(errInternalError)
	ErrCache         = errors.New(errCache)
	ErrAdapter       = errors.New(errAdapter)
	ErrNotFound      = errors.New(errNotFound)
)
