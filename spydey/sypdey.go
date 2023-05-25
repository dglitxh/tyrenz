package spydey

import (
	"io/fs"
)

type File struct {
	Info fs.FileInfo
	file fs.File
} 