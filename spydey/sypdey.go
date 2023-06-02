package spydey

import (
	"encoding/json"
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



func Find (filename, dirname string) error {
	var isFound bool
	if err := os.Chdir(dirname); err != nil {
		return err
	}
	filepath.WalkDir(dirname, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if entry.Name() == filename {
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



func Crawl (allow_hidden bool, dirname string) error{
	if err := os.Chdir(dirname); err != nil {
		return err
	}
	tree := make(map[string][]string)
	count := 0
	filepath.WalkDir(dirname, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		dirs := (strings.Split(path, "/"))
		dir := strings.Join(dirs[0:len(dirs)-1], "/")
		isHidden := string(path[0]) == "." && len(path)>1
		if !allow_hidden && isHidden {
			count++
		}else {
			if entry.IsDir()  {
				tree[path] = append(tree[path], tree[path]...)
			}else {
			if _, ok := tree[dir]; ok {
				tree[dir] = append(tree[dir], dirs[len(dirs)-1]) 
			}else {
				fmt.Println(dir)
				tree[dir] = append(tree[dir], tree[dir]...)
			}
			}
		}
		return nil
	})
	it, err := json.MarshalIndent(tree, " ", " "); if err != nil {
		return err
	}
	os.WriteFile("spydey.json", []byte(it), 0644)
	fmt.Printf("%d hidden files were exempted, use '-a' flag to allow hidden files\n", count)
	return nil
}