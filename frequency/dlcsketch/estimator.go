package dlcsketch

import (
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/frequency"
	"github.com/koykov/pbtk/frequency/cmsketch"
)

const flagDLC = 2

func NewEstimator[T pbtk.Hashable](config *cmsketch.Config) (frequency.Estimator[T], error) {
	if config == nil {
		return nil, pbtk.ErrInvalidConfig
	}
	config.WithFlag(flagDLC, true)
	return cmsketch.NewEstimator[T](config)
}
