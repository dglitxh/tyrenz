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


func Find (name string) error {
	filepath.WalkDir(".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if entry.Name() == name {
			fmt.Println(entry.Name(), "found @ ", Gwd()+"/"+path)
		}
		return nil
	})
	return nil
}

func Crawl (dir string) []File {
	var items []File
	
	return items
}