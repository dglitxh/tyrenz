package spydey

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Info fs.FileInfo
	Filename string
} 


func Gwd () string{
	wd, err := os.Getwd(); if err != nil {
		fmt.Println(err)
	}
	return wd
}


func Find (name string) error {
	var isFound bool
	filepath.WalkDir(".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if entry.Name() == name {
			fmt.Println(entry.Name(), "found @ ", Gwd()+"/"+path)
			isFound = true
			fmt.Println(strings.Split(path, "/"), path)
		}
		return nil
	})
	 if  !isFound {
		fmt.Println("File was not found in this directory")
		return errors.New("File not found")
	 }
	return nil
}

func Crawl (dir string) []File {
	var items []File
	
	return items
}