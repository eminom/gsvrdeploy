package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func doOneCopy(p interface{}) interface{} {
	ccmd := p.(*copyCmd)
	//fmt.Printf("Do-copy <%v> to <%v>\n", ccmd.abs_path, ccmd.newName)
	tPath := fmt.Sprintf("%v\\%v", conf_getTargetPath(), ccmd.newName)
	tPath = strings.Replace(tPath, "/", "\\", -1)
	//~ you may encrypt to the final destination in this phase.
	err := exec.Command("cmd",
		"/C",
		"copy",
		ccmd.abs_path,
		tPath,
		"/Y",
	).Run()
	if nil != err {
		fmt.Println("Error for copy")
		fmt.Printf("%v\n%v\n", ccmd.abs_path, ccmd.newName)
		fmt.Printf("to %v\n", tPath)
		panic(err)
	}
	return &copyCmdOut{}
}
