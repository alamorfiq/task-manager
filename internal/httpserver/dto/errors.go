package dto

import "errors"

var (
	ErrTitleRequired = errors.New("title is required")
	ErrInvalidID     = errors.New("invalid task id")
)

