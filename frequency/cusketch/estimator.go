package cusketch

import (
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/frequency"
	"github.com/koykov/pbtk/frequency/cmsketch"
)

const flagConservativeUpdate = 1<<iota - 1

func NewEstimator[T pbtk.Hashable](config *cmsketch.Config) (frequency.Estimator[T], error) {
	if config == nil {
		return nil, pbtk.ErrInvalidConfig
	}
	config.WithFlag(flagConservativeUpdate, true)
	return cmsketch.NewEstimator[T](config)
}
