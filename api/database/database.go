package database

import "errors"

var (
	ErrDuplicate = errors.New("record already exists")
)
