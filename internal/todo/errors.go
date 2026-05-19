package todo

import "errors"

var (
	ErrNotFound     = errors.New("task not found")
	ErrUnavailable  = errors.New("todo storage unavailable")
	ErrInvalidTitle = errors.New("task title is required")
)
