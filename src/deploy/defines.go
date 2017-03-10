
package main

type fileUno struct {
    abs_path string
    relpath string
}

type fileUnoOut struct {
    abs_path string
    relpath string
}

type calcPackIn struct {
    abs_path string
    relpath  string
}

type calcPackOut struct {
    hashName string
    size     int
    xxhash   string
    relpath  string
    abs_path  string
}

type copyCmd struct {
    newName  string
    abs_path string
}

type copyCmdOut struct {}