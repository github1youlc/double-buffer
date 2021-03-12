package dbuf

import (
	"io"
	"time"
)

// FileLoadSystem
type FileLoadSystem interface {
	// Detect the latest file and its modify time
	Detect(url string) (filePath string, modify time.Time, err error)

	// Read file
	Read(filePath string) (io.Reader, error)
}

// FileLoader
type FileLoader struct {
	url        string
	system     FileLoadSystem
	alloc      allocFunc
	load       loadFunc
	lastModify time.Time
}

// NewFileLoader
func NewFileLoader(system FileLoadSystem, url string, alloc allocFunc, load loadFunc) *FileLoader {
	return &FileLoader{
		url:    url,
		system: system,
		alloc:  alloc,
		load:   load,
	}
}

type allocFunc func() interface{}
type loadFunc func(reader io.Reader, i interface{}) error

// Load implement loader.Load
func (loader *FileLoader) Load(i interface{}) (bool, error) {
	filePath, modifyTime, err := loader.system.Detect(loader.url)
	if err != nil {
		return false, err
	}

	detectNew := modifyTime.After(loader.lastModify)

	if !detectNew {
		return false, nil
	}

	reader, err := loader.system.Read(filePath)
	if err != nil {
		return false, err
	}

	if err := loader.load(reader, i); err != nil {
		return false, err
	}

	return detectNew, nil
}

// Alloc implement loader.Alloc
func (loader *FileLoader) Alloc() interface{} {
	return loader.alloc()
}
