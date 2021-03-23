package dbuf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testLoader struct {
}

func (l *testLoader) Load(i interface{}) (bool, error) {
	return true, nil
}

func (l *testLoader) Alloc() interface{} {
	return 1
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
