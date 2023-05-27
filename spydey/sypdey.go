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
	
	return nil
}

func Crawl (dir string) []File {
	var items []File
	
	return items
}