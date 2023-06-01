package spydey

import (
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
			fmt.Println(entry.Name(), "found @", Gwd()+"/"+path)
			isFound = true
			fmt.Println(strings.Split(path, "/"), path)
		}
		return nil
	})
	 if  !isFound {
		fmt.Println("File was not found in this directory")
	 }
	return nil
}

func Crawl () []File {
	tree := make(map[string][]string)
	filepath.WalkDir(".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if entry.IsDir() {
			tree[entry.Name()] = append(tree[entry.Name()], "")
		}else {
			dir_arr := strings.Split(path, "/")
			if len(dir_arr) > 1 {
				if _, ok := tree[path];  ok {
					tree[path] = append(tree[path], path)
				}
			}
		}
		return nil
	})
	fmt.Println(tree)
	return nil
}