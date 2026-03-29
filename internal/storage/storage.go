package storage

import "errors"

var (
	ErrorAliasNotFound      = errors.New("alias not found")
	ErrorAliasAlreadyExists = errors.New("alias already exists")
)
