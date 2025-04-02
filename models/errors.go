package models

import "errors"

var (
	// A common pattern is to add the package as a prefix to the error for
	// context.
	ErrEmailTaken       = errors.New("models: email address is already in use")
	ErrPasswordMismatch = errors.New("models: password mismatch")
	ErrNotFound         = errors.New("models: resource not found")
)
