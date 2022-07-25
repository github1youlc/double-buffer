package dbuf

import (
	"io"
	"time"
)

// FileLoadSystem 文件加载系统接口
type FileLoadSystem interface {
	// Detect the latest file and its modify time
	Detect(url string) (filePath string, modify time.Time, err error)

	// Read file
	Read(filePath string) (io.Reader, error)
}

// FileLoader 文件加载器
type FileLoader struct {
	url        string
	system     FileLoadSystem
	load       LoadFunc
	lastModify time.Time
}

// NewFileLoader 创建文件加载器
func NewFileLoader(system FileLoadSystem, url string, load LoadFunc) *FileLoader {
	return &FileLoader{
		url:    url,
		system: system,
		load:   load,
	}
}

// LoadFunc load from reader
type LoadFunc func(reader io.Reader, i interface{}) error

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
	loader.lastModify = modifyTime

	return detectNew, nil
}
