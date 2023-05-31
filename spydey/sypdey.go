package spydey

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type File struct {
	Info fs.FileInfo
	filename string
} 


func Gwd () string{
	wd, err := os.Getwd(); if err != nil {
		fmt.Println(err)
	}
	
	return wd
}


func Find (name string) error {
	var mfile File
	var isFound bool
	err := filepath.WalkDir(".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		info, err := entry.Info(); if err != nil {
			return err
		}
		if entry.Name() == name {
			fmt.Println(entry.Name(), "found @ ", Gwd()+"/"+path)
			mfile.Info = info
			mfile.filename  = entry.Name()
			isFound = true
		}
		return nil
	})
	fmt.Println(err)
	 if  !isFound {
		fmt.Println("File was not found in this directory")
		return err
	 }
	return nil
}

func Crawl (dir string) []File {
	var items []File
	
	return items
}