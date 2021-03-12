package dbuf

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileLoader(t *testing.T) {
	filePath := path.Join(t.TempDir(), t.Name())

	f, err := os.Create(filePath)
	assert.NoError(t, err)
	_, _ = f.Write([]byte("abc"))
	assert.NoError(t, f.Close())

	fl := NewFileLoader(LocalSingleFileLoadSystem, filePath, func() interface{} {
		var s string
		return &s
	}, func(reader io.Reader, i interface{}) error {
		content, _ := ioutil.ReadAll(reader)
		s := i.(*string)
		*s = string(content)
		return nil
	})

	buffer := NewDoubleBuffer(fl)
	assert.True(t, buffer.load())
	assert.EqualValues(t, "abc", *buffer.Data().(*string))

	f, err = os.Create(filePath)
	assert.NoError(t, err)
	_, _ = f.Write([]byte("def"))
	assert.NoError(t, f.Close())

	assert.True(t, buffer.load())
	assert.EqualValues(t, "def", *buffer.Data().(*string))
}