package hyperloglog

import "errors"

var ErrInvalidPrecision = errors.New("precision must be in range [4..18]")
