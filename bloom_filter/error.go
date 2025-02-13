package bloom

import "errors"

var (
	ErrBadConfig = errors.New("invalid or empty config")
	ErrNoHasher  = errors.New("no hasher provided")
)
