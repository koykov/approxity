package lsh

import "errors"

var (
	ErrNoShingler = errors.New("no shingler provided")
	ErrZeroK      = errors.New("zero K provided")
	ErrZeroB      = errors.New("zero B provided")
	ErrBigB       = errors.New("too big B provided, must be less than 64")
)
