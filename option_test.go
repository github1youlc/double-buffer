package dbuf

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOption(t *testing.T) {
	opt := newOption(
		WithErrCallback(func(err error) {}),
		WithReloadInterval(time.Second),
		WithInitCallback(func() {}),
	)

	assert.EqualValues(t, time.Second, opt.reloadInterval)
	assert.Len(t, opt.errCallback, 1)
	assert.NotNil(t, opt.initCallback)
}
