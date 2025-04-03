package tinylfu

import "errors"

var ErrDecayRange = errors.New("decay factor or soft factor must be in range (0..1)")
