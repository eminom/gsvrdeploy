package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func doOneCheck(p interface{}) interface{} {
	ptr := p.(*fileUno)
	fi, err := os.Stat(ptr.abs_path)
	if nil != err {
		panic(err)
	}
	ext := filepath.Ext(fi.Name())
	if `.lua` == strings.ToLower(ext) {
		pn := strings.TrimSuffix(fi.Name(), ext)
		if !isNameExcepted(pn) {
			doXXTeaToUno(ptr.abs_path, ptr.abs_path)
		} else {
			fmt.Printf("Except for <%v>\n", pn)
		}
	}
	return &fileUnoOut{ptr.abs_path, ptr.relpath}
}
