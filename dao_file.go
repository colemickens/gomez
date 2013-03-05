package main

import (
	"github.com/eaigner/hood"
)

type GenericDao struct {
	// inject hood hd or just use the global one?
}

func (gd GenericDao) AddNewFile(f *FileRecord) error {
	_, err := hd.Save(f)
	if err != nil {
		panic(err)
	}
	return nil
}

func (gd GenericDao) PathExists(path string) bool {
	var files []FileRecord
	err := hd.Where("path", "=", path).Limit(1).Find(&files)
	if err != nil {
		panic(err)
	}
	return (len(files) > 0)
}
