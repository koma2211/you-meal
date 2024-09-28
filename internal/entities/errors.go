package entities

import "errors"

var (
	ErrEmptyBurgers   = errors.New("burgers are empty")
	ErrImageNotExists = errors.New("image not exists")
)
