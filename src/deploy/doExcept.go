package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	excepted = make(map[string]bool)
)

func init() {
	path := ".excepted"
	fin, err := os.Open(path)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "Cannot open %s", path)
		//panic(err)
		fmt.Fprintf(os.Stderr, "You may need a exception list")
		return
	}
	defer fin.Close()

	rd := bufio.NewReader(fin)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		//~ Remember to trim the \r and \n.(especially \r)
		line = strings.Trim(line, "\t \n\r")
		excepted[line] = true
		//fmt.Printf("in:%s\n", line)
	}
}

func isNameExcepted(name string) bool {
	_, ok := excepted[name]
	return ok
}
