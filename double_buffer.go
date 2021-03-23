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
}

// DoubleBuffer double buffer, used to hot load
type DoubleBuffer struct {
	loader      Loader
	bufferData  []interface{}
	curIndex    int32
	errCallback []func(error)
	started     int32
	opt         *option

	initOnce sync.Once
}

// Alloc allocate the target object
type Alloc func() interface{}

// NewDoubleBuffer create double buffer object
func NewDoubleBuffer(loader Loader, alloc Alloc, option ...Option) *DoubleBuffer {
	opt := newOption(option...)

	b := &DoubleBuffer{
		loader:   loader,
		curIndex: 0,
		started:  0,
		opt:      opt,
	}
	b.bufferData = append(b.bufferData, alloc(), alloc())
	return b
}

// load load data
func (b *DoubleBuffer) Start() {
	b.load()
	go func() {
		for {
			time.Sleep(b.opt.reloadInterval)
			_ = b.load()
		}
	}()

	return
}

func (b *DoubleBuffer) load() bool {
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
func (b *DoubleBuffer) Data() interface{} {
	ci := atomic.LoadInt32(&b.curIndex)
	return b.bufferData[ci]
}
