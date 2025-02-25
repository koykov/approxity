package quotient

import "sync"

type Filter struct {
	conf *Config
	once sync.Once
	// todo implement me
}

func (f *Filter) Set(key any) error {
	// todo implement me
	return nil
}

func (f *Filter) HSet(hkey uint64) error {
	// todo implement me
	return nil
}

func (f *Filter) Unset(key any) error {
	// todo implement me
	return nil
}

func (f *Filter) HUnset(hkey uint64) error {
	// todo implement me
	return nil
}

func (f *Filter) Contains(key any) bool {
	// todo implement me
	return false
}

func (f *Filter) HContains(hkey uint64) bool {
	// todo implement me
	return false
}

func (f *Filter) Reset() {
	// todo implement me
}
