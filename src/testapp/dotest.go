package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func getTargetVersionFile() string {
	return "H:/gevents/hashv/version_x.txt"
}

func getTargetPath() string {
	return "H:/gevents/resfolder/KillAll"
}

func ensureDir(t string) {
	toMk := filepath.Dir(t)
	err := os.MkdirAll(toMk, os.ModeDir)
	if nil != err {
		panic(err)
	}
	fmt.Printf("%v is made\n", toMk)
}

func main() {

	t := getTargetVersionFile()
	ensureDir(t)
	ensureDir(getTargetPath())
}
