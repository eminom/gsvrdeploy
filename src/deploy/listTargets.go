package main

import (
	"os"
	"path/filepath"
	"strings"
)

/////////////////////////////////////////////////
//~ The list for current deploying version
func getAllFiles(root string) []*fileUno {
	files := make([]*fileUno, 0)
	file_count, dir_count := 0, 0
	err := filepath.Walk(root, func(patho string, fi os.FileInfo, inpErr error) (err error) {
		if nil != inpErr {
			panic(inpErr)
		}
		path := strings.TrimPrefix(patho, root)
		//~ both
		path = strings.TrimPrefix(path, "/")
		path = strings.TrimPrefix(path, "\\")
		if fi.IsDir() {
			dir_count += 1
		} else {
			file_count += 1
			abs, _ := filepath.Abs(filepath.Join(root, path))
			files = append(files, &fileUno{abs, path})
		}
		return nil
	})
	if nil != err {
		panic(err)
	}
	return files
}
