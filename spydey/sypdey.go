package spydey

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type File struct {
	Info fs.FileInfo
	file fs.File
} 


func Gwd () string{
	wd, err := os.Getwd(); if err != nil {
		fmt.Println(err)
	}
	
	return wd
}


func Search (item string) error {
	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}

func Crawl (dir string) []File {
	var items []File
	
	return items
}