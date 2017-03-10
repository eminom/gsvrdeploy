

package main

import (
    "fmt"
    "os/exec"
)

func doOneCopy(p interface{}) interface{} {
    ccmd := p.(*copyCmd)
    //fmt.Printf("Do-copy <%v> to <%v>\n", ccmd.abs_path, ccmd.newName)
    tPath := fmt.Sprintf("%v\\%v", getTargetPath(), ccmd.newName)
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
        panic(err)
    }
    return &copyCmdOut{}
}