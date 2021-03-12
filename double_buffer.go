package dbuf

import (
	"sync"
	"sync/atomic"
	"time"
)

// Loader
type Loader interface {
	// Load data to object
	// If data is updated, return updated should be true
	Load(i interface{}) (updated bool, err error)

	// Alloc allocate the target object
	Alloc() interface{}
}

// doubleBuffer double buffer, used to hot load
type doubleBuffer struct {
	loader      Loader
	bufferData  []interface{}
	curIndex    int32
	errCallback []func(error)
	started     int32
	opt         *option

	initOnce sync.Once
}

// NewDoubleBuffer create double buffer object
func NewDoubleBuffer(loader Loader, option ...Option) *doubleBuffer {
	opt := newOption(option...)

	b := &doubleBuffer{
		loader:   loader,
		curIndex: 0,
		started:  0,
		opt:      opt,
	}
	b.bufferData = append(b.bufferData, loader.Alloc(), loader.Alloc())
	return b
}

// load load data
func (b *doubleBuffer) Start() {
	b.load()
	go func() {
		for {
			time.Sleep(b.opt.reloadInterval)
			_ = b.load()
		}
	}()

	return
}

func (b *doubleBuffer) load() bool {
	ci := 1 - atomic.LoadInt32(&b.curIndex)
	updated, err := b.loader.Load(b.bufferData[ci])
	if err != nil {
		for _, cb := range b.opt.errCallback {
			cb(err)
		}
		return false
	}
	if updated {
		atomic.StoreInt32(&b.curIndex, ci)
		b.initOnce.Do(
			func() {
				if b.opt.initCallback != nil {
					b.opt.initCallback()
				}
			})
	}

	return updated
}

// Data get latest data
func (b *doubleBuffer) Data() interface{} {
	ci := atomic.LoadInt32(&b.curIndex)
	return b.bufferData[ci]
}
