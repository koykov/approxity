package lossy

import "errors"

var (
	ErrZeroEpsilon = errors.New("epsilon must be greater than zero")
	ErrBadEpsilon  = errors.New("epsilon must be greater that support")
)
