package store

import "errors"

var (
	ErrKeyNotFound = errors.New("key does not exist in store")
	ErrKeyExpired = errors.New("key has expired")
)

