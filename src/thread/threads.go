package threads

import (
	"fmt"
	"reflect"
	"sync"
	_ "time"
)

type WorkRoutine func(interface{}) interface{}

type Processor struct {
	wr WorkRoutine
}

var processors map[interface{}]*Processor

func DoInit() {
	processors = make(map[interface{}]*Processor)
}

func RegisterWorkRoutine(m interface{}, wr WorkRoutine) {
	typeEste := reflect.TypeOf(m)
	processors[typeEste] = &Processor{wr}
}

func getProcessor(m interface{}) *Processor {
	return processors[reflect.TypeOf(m)]
}

type workUnit struct {
	Intv       int
	RefCounter int
	Name       string
	IsClosing  *chan bool
	wg         sync.WaitGroup
}

func makeWorkUnit(n int) *workUnit {
	rv := &workUnit{Intv: n, Name: fmt.Sprintf("<%v>", n)}
	isClosin := make(chan bool)
	rv.IsClosing = &isClosin
	return rv
}

func (self *workUnit) StartWork(in chan interface{}, out chan interface{}) {
	self.wg.Add(1)
	go self._workRoutine(in, out)
}

func (self *workUnit) _workRoutine(in chan interface{}, out chan interface{}) {
	refInt := &self.RefCounter
Main100:
	for {
		select {
		case uno := <-in:
			res := callToFunc(uno)
			out <- res
			*refInt++
		case <-*self.IsClosing:
			break Main100
		}
	}
	self.wg.Done()
}

func callToFunc(m interface{}) interface{} {
	proc := getProcessor(m)
	return proc.wr(m)
}

type WorkUnitPool struct {
	objs []*workUnit
}

func MakeWorkPool() *WorkUnitPool {
	rv := &WorkUnitPool{}
	return rv
}

func (self *WorkUnitPool) StartPool(size int, in chan interface{}, out chan interface{}) {
	self.objs = make([]*workUnit, size)
	wu := self.objs
	for i := 0; i < len(wu); i++ {
		wu[i] = makeWorkUnit(i)
	}
	for _, obj := range wu {
		obj.StartWork(in, out)
	}
}

func (self *WorkUnitPool) Shutdown() {
	for _, obj := range self.objs {
		*obj.IsClosing <- true
		//fmt.Printf("Waiting for <%v>\n", obj.Name)
		obj.wg.Wait()
	}
	//? fmt.Println("Pool closed")
}

func (self *WorkUnitPool) PrintInfo() {
	for _, obj := range self.objs {
		fmt.Printf("<%v> is waken for %v time(s)\n", obj.Name, obj.RefCounter)
	}
}

const (
	DefaultWorkCount = 256
)
