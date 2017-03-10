package main

import (
	"os"
	"os/signal"
	"thread"
	"time"
	"fmt"
)


func runEvaluate(taskName string, proc func()){
	start := time.Now()
	proc()
	fmt.Printf("%v runs for %v\n", taskName, time.Since(start))
}

////////////////////////////////////////////////////////////////////////////
func masterTask() {
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

	//~ all in
	for _, f := range allFiles {
		inputChan <- f
	}
	left := expected * 3 + 1

Main100:
	for left > 0 {
		select {
		case esto := <-outs:
			var dc int = 1
			switch esto.(type) {
			case *copyCmdOut:
				break
			case *fileUnoOut:
				eso := esto.(*fileUnoOut)
				inputChan <- &calcPackIn{abs_path:eso.abs_path, relpath:eso.relpath}
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
	fmt.Println("done")
}


func main() {
	runEvaluate("deploying", masterTask)
}
