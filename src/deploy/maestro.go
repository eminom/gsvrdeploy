
package main

import (
	"fmt"
	"os/signal"
	"os"
	"thread"
    "path/filepath"
    "crypto/md5"
    "io"
    "strings"
    "github.com/OneOfOne/xxhash"
)


func CalcFileMd5(path string)string{
    fin, err := os.Open(path)
    if nil != err {
        panic(err)
    }
    defer fin.Close()
    h := md5.New()
    _, e2 := io.Copy(h, fin);
    if nil != e2 {
        panic(err)
    }
    // 总是32位宽的16进制的数字;
    return fmt.Sprintf("%x", h.Sum(nil))
}

func CalcFileXXHash(path string)string{
    fin, err := os.Open(path)
    if nil != err {
        panic(err)
    }
    defer fin.Close()
    h := xxhash.NewS32(10241024)   // New32(注意是10进制)
    _, e2 := io.Copy(h, fin)
    if nil != e2 {
        panic(e2)
    }
    return fmt.Sprintf("%08x", h.Sum32())
}

//定制标准ENTRY
func CalcFileEnt(abs_path string, relpath string, fi os.FileInfo)string{
    md5 := CalcFileMd5(abs_path)
    xxhash := CalcFileXXHash(abs_path)
    relpath = strings.Replace(relpath, "\\", "/", -1)
    return fmt.Sprintf("%s%s\t%d\t%s\t%s", md5, filepath.Ext(fi.Name()), fi.Size(), xxhash, relpath)
    //return fmt.Sprintf("%s", xxhash)
}

type CalcPack struct {
	abs_path string
	relpath string
	fi os.FileInfo  // interface
}

func (self *CalcPack)Print(){
	fmt.Printf("<%s>\n", self.fi.Name())
}

func hacerLaCama(p interface{})interface{}{
	fpath := p.(*CalcPack)
	res := CalcFileEnt(fpath.abs_path, fpath.relpath, fpath.fi)
	return &res
}

func getAllFiles(root string)[]*CalcPack{
	files := make([]*CalcPack, 0)
    file_count, dir_count := 0, 0
    err := filepath.Walk(root, func(patho string, fi os.FileInfo, inpErr error)(err error){
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
            files = append(files, &CalcPack{abs, path, fi})
        }
        return nil
    })
    if nil != err {
        panic(err)
    }
    return files
}

////////////////////////////////////////////////////////////////////////////
func main() {
    if len(os.Args) < 2 {
        fmt.Println("Need more input parameter <folder-name>")
        os.Exit(1)
    }
    target_dir := os.Args[1]
    di, dErr := os.Stat(target_dir)
    if nil != dErr {
        panic(dErr)
    }
    if ! di.IsDir() {
        fmt.Printf("`%s' is not a valid directory.\n", target_dir)
        os.Exit(2)
    }
	allFiles := getAllFiles(target_dir)
	expected := len(allFiles)
	resList := make(chan interface{}, expected)
	workpool := threads.MakeWorkPool()
	wakeUno := workpool.Start(threads.DefaultWorkCount, expected, hacerLaCama, resList)
	fullOuts := make([]string, 0)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

    //~ 全部放进去再说
	for _, f := range allFiles{
		wakeUno <- f
	}
	left := expected
	Main100: for ;left>0;{
		select {
		case resOut := <- resList:
			left--
			fullOuts = append(fullOuts, *resOut.(*string))
		case <- c:
			break Main100
		default: //~ this is so necessary.(or blocked above)
		}
	}
	workpool.Shutdown()
    //workpool.PrintInfo()
	for _, v := range fullOuts{
        fmt.Println(v)
    }

}