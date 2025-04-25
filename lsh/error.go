package lsh

import "errors"

var (
	ErrNoShingler = errors.New("no shingler provided")
	ErrZeroK      = errors.New("zero K provided")
)
