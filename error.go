package bloom

import "errors"

var (
	ErrBadConfig = errors.New("invalid or empty config")
	ErrBadPolicy = errors.New("unsupported policy provided")
	ErrSetAccess = errors.New("cannot set access due to policy")
)
