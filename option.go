package dbuf

import "time"

type option struct {
	errCallback    []func(error)
	reloadInterval time.Duration
	initCallback   func()
}

var DefaultReloadInterval = time.Second * 5

func newOption(modify ...Option) *option {
	opt := &option{
		reloadInterval: DefaultReloadInterval,
	}

	for _, m := range modify {
		m(opt)
	}

	return opt
}

// Option is used to modify default option
type Option func(option *option)

// WithErrCallback set error callback
func WithErrCallback(errCallback ...func(error)) Option {
	return func(option *option) {
		option.errCallback = errCallback
	}
}

// WithReloadInterval set reload interval
func WithReloadInterval(interval time.Duration) Option {
	return func(option *option) {
		option.reloadInterval = interval
	}
}

// WithInitCallback set initCallback
// If the call should wait double buffer initialization (at least load data once) done
//  done := make(chan struct{}, 1)
//	buffer := NewDoubleBuffer(&testLoader{}, WithInitCallback(
//		func() {
//			done <- struct{}{}
//		}))
//
//	buffer.Start()
//	<-done
func WithInitCallback(initCallback func()) Option {
	return func(option *option) {
		option.initCallback = initCallback
	}
}
