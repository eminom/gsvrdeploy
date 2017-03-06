package main

import (
	"crypto/md5"
	"fmt"
	"github.com/OneOfOne/xxhash"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"thread"
	"time"
)

type calcPackIn struct {
	abs_path string
	relpath  string
	fi       os.FileInfo // interface
}

type calcPackOut struct {
	hashName string
	size     int
	xxhash   string
	relpath  string
	abspath  string
}

func (self *calcPackOut) toStr() string {
	return fmt.Sprintf("%v\t%v\t%v\t%v", self.hashName, self.size, self.xxhash, self.relpath)
}

func CalcFileMd5(path string) string {
	fin, err := os.Open(path)
	if nil != err {
		panic(err)
	}
	defer fin.Close()
	h := md5.New()
	_, e2 := io.Copy(h, fin)
	if nil != e2 {
		panic(err)
	}
	// MD5:总是32位宽的16进制的数字;
	return fmt.Sprintf("%x", h.Sum(nil))
}

func CalcFileXXHash(path string) string {
	fin, err := os.Open(path)
	if nil != err {
		panic(err)
	}
	defer fin.Close()
	h := xxhash.NewS32(10241024) // New32(注意是10进制)
	_, e2 := io.Copy(h, fin)
	if nil != e2 {
		panic(e2)
	}
	return fmt.Sprintf("%08x", h.Sum32())
}

//定制标准ENTRY
func calcFileEnt(abs_path string, relpath string, fi os.FileInfo) interface{} {
	md5 := CalcFileMd5(abs_path)
	xxhash := CalcFileXXHash(abs_path)
	relpath = strings.Replace(relpath, "\\", "/", -1)
	//return fmt.Sprintf("%s%s\t%d\t%s\t%s", md5, filepath.Ext(fi.Name()), fi.Size(), xxhash, relpath)
	outName := fmt.Sprintf("%s%s", md5, filepath.Ext(fi.Name()))
	return &calcPackOut{outName, int(fi.Size()), xxhash, relpath, abs_path}
}

func hacerLaCama(p interface{}) interface{} {
	fpath := p.(*calcPackIn)
	return calcFileEnt(fpath.abs_path, fpath.relpath, fpath.fi)
}

var meTargetPath string

func formatTargetPath(n string) string {
	return fmt.Sprintf("%v\\%v", meTargetPath, n)
}

type copyCmd struct {
	newName  string
	abs_path string
}

func (self *copyCmd) do() {
	err := exec.Command("cmd",
		"/C",
		"copy",
		self.abs_path,
		formatTargetPath(self.newName),
		"/Y",
	).Run()
	if nil != err {
		fmt.Println("Error for copy")
		fmt.Printf("%v\n%v\n", self.abs_path, self.newName)
		panic(err)
	}
}

func doOneCopy(p interface{}) interface{} {
	ccmd := p.(*copyCmd)
	//fmt.Printf("Do-copy <%v> to <%v>\n", ccmd.abs_path, ccmd.newName)
	ccmd.do()
	return 1
}

/////////////////////////////////////////////////
//~ The list for current deploying version
func getAllFiles(root string) []*calcPackIn {
	files := make([]*calcPackIn, 0)
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
			files = append(files, &calcPackIn{abs, path, fi})
		}
		return nil
	})
	if nil != err {
		panic(err)
	}
	return files
}

func doWriteToVersionFile(arr []string) {
	name := getTargetVersionFile()
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	for _, line := range arr {
		fmt.Fprintf(f, "%v\n", line)
	}
}

////////////////////////////////////////////////////////////////////////////
func main() {
	start := time.Now()
	if len(os.Args) < 2 {
		fmt.Println("Need more input parameter <folder-name>")
		os.Exit(1)
	}
	target_dir := os.Args[1]
	di, dErr := os.Stat(target_dir)
	if nil != dErr {
		panic(dErr)
	}
	if !di.IsDir() {
		fmt.Printf("`%s' is not a valid directory.\n", target_dir)
		os.Exit(2)
	}
	allFiles := getAllFiles(target_dir)
	expected := len(allFiles)

	meTargetPath = getTargetPath()
	//~  THE OUPUT
	outs := make(chan interface{}, 12)

	//~ register methods for route's processing
	threads.DoInit()
	threads.RegisterWorkRoutine(&calcPackIn{}, hacerLaCama)
	threads.RegisterWorkRoutine(&copyCmd{}, doOneCopy)

	inputChan := make(chan interface{}, expected*2)
	workpool := threads.MakeWorkPool()
	workpool.StartPool(threads.DefaultWorkCount, inputChan, outs)

	fullOuts := make([]string, 0)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	//~ all in
	for _, f := range allFiles {
		inputChan <- f
	}
	left := expected * 2

Main100:
	for left > 0 {
		select {
		case resOut := <-outs:
			var dc int = 1
			switch resOut.(type) {
			case int:
				break
			case *calcPackOut:
				po := resOut.(*calcPackOut)
				fullOuts = append(fullOuts, po.toStr())
				inputChan <- &copyCmd{po.hashName, po.abspath}
				break
			default:
				panic("error")
				dc = 0
			}
			left -= dc
		case <-c:
			break Main100
		default: //~ this is so necessary.(or blocked above)
		}
	}
	workpool.Shutdown()
	//workpool.PrintInfo()
	doWriteToVersionFile(fullOuts)
	fmt.Println(time.Since(start))
}
