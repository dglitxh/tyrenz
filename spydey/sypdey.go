package spydey

import (
	"io/fs"
)

type File struct {
	Info fs.FileInfo
	file fs.File
} 



func Search (item string) error {

	return nil
}

func Crawl (dir string) []File {
	var items []File
	
	return items
}