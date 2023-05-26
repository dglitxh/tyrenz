package spydey

import (
	"io/fs"
)

type File struct {
	Info fs.FileInfo
	file fs.File
} 

type MyFiles []File

func (f *MyFiles) Search (item string) error {

	return nil
}