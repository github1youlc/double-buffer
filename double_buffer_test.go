package dbuf

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testLoader struct {
}

func (l *testLoader) Load(i interface{}) (bool, error) {
	return true, nil
}

func TestFileDoubleBuffer(t *testing.T) {
	done := make(chan struct{}, 1)
	buffer := NewDoubleBuffer(
		&testLoader{},
		func() interface{} { return 1 },
		WithInitCallback(
			func() {
				done <- struct{}{}
			}))

	buffer.Start()
	<-done
	assert.EqualValues(t, 1, buffer.Data())
}

type updatedLoader struct {
	count int
}

func (l *updatedLoader) Load(i interface{}) (bool, error) {
	l.count++
	d := i.(*int)
	*d = l.count
	return true, nil
}

func TestFileDoubleBufferUpdateCallback(t *testing.T) {
	done := make(chan struct{}, 1)
	updatedData := make(chan interface{}, 1)
	var buffer *DoubleBuffer
	buffer = NewDoubleBuffer(
		&updatedLoader{},
		func() interface{} { return new(int) },
		WithInitCallback(
			func() {
				done <- struct{}{}
			}),
		WithReloadInterval(time.Millisecond),
		WithUpdatedCallback(
			func() {
				updatedData <- buffer.Data()
			}),
	)

	buffer.Start()
	<-done
	for i := 1; i < 10; i++ {
		a := <-updatedData
		assert.EqualValues(t, i, *(a.(*int)))
	}
}
