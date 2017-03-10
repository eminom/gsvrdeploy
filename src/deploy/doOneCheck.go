

package main

import (
    _ "fmt"
    "os"
    "path/filepath"
    "strings"
)

func doOneCheck(p interface{})interface{}{
    ptr := p.(*fileUno)
    fi, err := os.Stat(ptr.abs_path)
    if nil != err {
        panic(err)
    }
    ext := filepath.Ext(fi.Name())
    if `.lua` == strings.ToLower(ext) {
        doXXTeaToUno(ptr.abs_path, ptr.abs_path)
    }
    return &fileUnoOut{ptr.abs_path, ptr.relpath}
}