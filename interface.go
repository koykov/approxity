package amq

type Interface interface {
	Set(key any) error
	Unset(key any) error
	Contains(key any) bool
	Reset()
}
