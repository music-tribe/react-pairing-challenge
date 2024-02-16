package database

import "errors"

var (
	ErrDuplicate = errors.New("record already exists")
	ErrNotFound  = errors.New("record not found")
)
