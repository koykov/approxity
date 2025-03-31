package tinylfu

type ForceDecayNotifier interface {
	Notify() <-chan struct{}
}

type dummyForceDecayNotifier struct{}

func (dummyForceDecayNotifier) Notify() <-chan struct{} {
	return nil
}
