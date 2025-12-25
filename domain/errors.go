package domain

import "errors"

var (
	ErrNotFound             = errors.New("wallet not found")
	ErrInternal             = errors.New("internal server error")
	ErrInvalidID            = errors.New("invalid id")
	ErrInvalidBalance       = errors.New("invalid balance")
	ErrInvalidOperation     = errors.New("invalid operation")
	ErrInsufficientBalance  = errors.New("insufficient balance")
	ErrInvalidOperationType = errors.New("invalid operation type")
)
