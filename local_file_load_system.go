package dbuf

import (
	"io"
	"os"
	"time"
)

var (
	LocalSingleFileLoadSystem = &localSingleFileLoadSystem{}
)

type localSingleFileLoadSystem struct {

}

func (sys *localSingleFileLoadSystem) Detect(url string) (filePath string, modify time.Time, err error) {
	var statInfo os.FileInfo
	statInfo, err = os.Stat(url)
	if err != nil {
		return
	}

	filePath = url
	modify = statInfo.ModTime()
	return
}

func (sys *localSingleFileLoadSystem) Read(filePath string) (io.Reader, error) {
	return os.Open(filePath)
}
