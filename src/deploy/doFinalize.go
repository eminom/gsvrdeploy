package main

import (
	"fmt"
	"os"
)

func doWriteToVersionFile(arr []string) {
	f, err := os.OpenFile(conf_getTargetVersionFile(), os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	for _, line := range arr {
		fmt.Fprintf(f, "%v\n", line)
	}
}
