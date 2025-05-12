package spacesaving

import "github.com/koykov/pbtk"

type hitter[T pbtk.Hashable] struct {
	conf *Config
}
