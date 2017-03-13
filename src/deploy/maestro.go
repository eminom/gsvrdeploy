package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"regexp"
	"thread"
	"time"
)

var version string
var input_folder string
var cdn_ip string

const (
	default_version    = "1.0.0"
	//deploy_target_base = "H:/GDWorks/test/g/v1/gevents/"
	deploy_target_base = `I:\JokerRush\AppX\gevents\`
)

func conf_getTargetVersionFile() string {
	return deploy_target_base + "resfolder/hashv/" + version + ".txt"
}

func conf_getTargetPath() string {
	return deploy_target_base + "resfolder/res"
}

func conf_getDistFile() string {
	return deploy_target_base + "resfolder/version.txt"
}

func init() {
	//fmt.Println("Inside main's init")
	flag.StringVar(&version, "version", "1.0.0", "distribution version string")
	flag.StringVar(&input_folder, "folder", "spinematch", "folder for distribution content")
	flag.Parse()
	if matched, _ := regexp.MatchString(`^\d+\.\d+\.\d+$`, version); !matched {
		version = default_version
		fmt.Printf("Using default version:%v\n", version)
	} else {
		fmt.Printf("version:%v\n", version)
	}
	fmt.Printf("input-folder:%v\n", input_folder)
}

type Distr struct {
	Version string `json:"version"`
	Cdn     string `json:"cdn"`
	Size    int    `json:"size"`
	Basever string `json:"basever"`
}

/* Method 1:
func (this Distr) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"version": this.version,
		"cdn":     this.cdn,
		"size":    this.size,
		"basever": this.basever,
	})
}
*/

func doWriteMyDistr() {
	d := Distr{
		version,
		"192.168.18.1:7991",
		1048,
		"1.0.X",
	}
	//~ OpenFile:  O_CREATE and O_WRONLY does not truncate this file.
	fout, err := os.Create(conf_getDistFile())
	if nil != err {
		panic(err)
	}
	defer fout.Close()
	b, e1 := json.Marshal(d)
	if nil != e1 {
		panic(e1)
	}
	fout.Write(b)
}

func runEvaluate(taskName string, proc func()) {
	start := time.Now()
	//defer fmt.Printf("%v runs for %v\n", taskName, time.Since(start))
	proc()
	fmt.Printf("%v runs for %v\n", taskName, time.Since(start))
}

type TaskCounter struct {
	counter map[reflect.Type]int
	allc    int
}

func newTaskCounter() *TaskCounter {
	rv := &TaskCounter{}
	rv.counter = make(map[reflect.Type]int)
	return rv
}

func (self *TaskCounter) incType(t reflect.Type) {
	if 0 == self.counter[t] {
		self.allc++
	}
	self.counter[t] += 1
}

func (self *TaskCounter) decType(t reflect.Type) {
	self.counter[t] -= 1
	if 0 == self.counter[t] {
		self.allc--
	}
}

func (self *TaskCounter) isDone() bool {
	return 0 == self.allc
}

////////////////////////////////////////////////////////////////////////////
func masterTask() {
	/*
		if len(os.Args) < 2 {
			fmt.Println("Need more input parameter <folder-name>")
			os.Exit(1)
		}
	*/

	//~ make target's folder for them.
	for _, path := range []string{conf_getTargetPath(), filepath.Dir(conf_getTargetVersionFile())} {
		mkErr := os.MkdirAll(path, os.ModeDir)
		if nil != mkErr {
			panic(mkErr)
		}
	}

	src_dir := input_folder
	di, dErr := os.Stat(src_dir)
	if nil != dErr {
		//panic(dErr)
		fmt.Fprintf(os.Stderr, "%v is not a directory.\n", src_dir)
		os.Exit(1)
	}
	if !di.IsDir() {
		fmt.Printf("`%s' is not a valid directory.\n", src_dir)
		os.Exit(2)
	}
	allFiles := getAllFiles(src_dir)
	expected := len(allFiles)

	//~  THE OUPUT
	outs := make(chan interface{}, 12)

	//~ register methods for route's processing
	threads.DoInit()
	threads.RegisterWorkRoutine(&calcPackIn{}, hacerLaCama)
	threads.RegisterWorkRoutine(&copyCmd{}, doOneCopy)
	threads.RegisterWorkRoutine(&fileUno{}, doOneCheck)

	inputChan := make(chan interface{}, expected*2)
	workpool := threads.MakeWorkPool()
	workpool.StartPool(threads.DefaultWorkCount, inputChan, outs)

	fullOuts := make([]string, 0)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	//counterMap := newTaskCounter()
	//~ all in
	for _, f := range allFiles {
		inputChan <- f
	}
	left := expected * 3
Main100:
	//for !counterMap.isDone() {
	for left > 0 {
		select {
		case esto := <-outs:
			var dc int = 1
			switch esto.(type) {
			case *copyCmdOut:
				break
			case *fileUnoOut:
				eso := esto.(*fileUnoOut)
				inputChan <- &calcPackIn{abs_path: eso.abs_path, relpath: eso.relpath}
				break
			case *calcPackOut:
				eso := esto.(*calcPackOut)
				fullOuts = append(fullOuts, eso.toStr())
				inputChan <- &copyCmd{eso.hashName, eso.abs_path}
				break
			default:
				panic("error")
				dc = 0
			}
			left -= dc
		case <-c:
			fmt.Println("User breaks")
			break Main100
		default: //~ this is so necessary.(or blocked above)
		}
	}
	workpool.Shutdown()
	//workpool.PrintInfo()
	doWriteToVersionFile(fullOuts)
	doWriteMyDistr()
	fmt.Println("done")
}

func main() {
	runEvaluate("deploying", masterTask)
}
