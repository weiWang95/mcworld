package util

import (
	"os"
	"path/filepath"
)

func AbsPath(dirDataName string) string {
	Mkdir(dirDataName)

	dirData, err := filepath.Abs(dirDataName)
	if err != nil {
		panic(err)
	}
	return dirData
}

func Mkdir(dirName string) {
	if err := os.MkdirAll(dirName, 0777); err != nil && !os.IsExist(err) {
		panic(err)
	}
}
