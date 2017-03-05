
package threads

import (
    "sync"
    _ "time"
    "fmt"
)

type workUnit struct {
    Intv int
    RefCounter int
    Name string
    IsClosing *chan bool
    WorkFunc interface{}
    wg sync.WaitGroup
}

func makeWorkUnit(n int) *workUnit{
    rv := &workUnit{Intv:n, Name:fmt.Sprintf("<%v>", n)}
    isClosin := make(chan bool)
    rv.IsClosing = &isClosin
    return rv
}

func (self *workUnit)Start(do func(interface{})interface{}, in chan interface{}, out chan interface{}) {
    self.WorkFunc = do
    self.wg.Add(1)
    go self._workRoutine(in, out)
}

func (self *workUnit)_workRoutine(in chan interface{}, out chan interface{}) {
    //intv := self.Intv
    refInt := &self.RefCounter
    worker := self.WorkFunc
    Main100:for {
        //time.Sleep(time.Duration(intv * 666) * time.Millisecond)
        select {
        case verUno := <- in:
            resI := worker.(func(interface{})interface{})(verUno)
            out <- resI
            *refInt++
        case <- *self.IsClosing:
            break Main100
        }
    }
    self.wg.Done()
}

type WorkUnitPool struct{
    objs []*workUnit
}

func MakeWorkPool() *WorkUnitPool {
    rv := &WorkUnitPool{}
    return rv
}

func (self *WorkUnitPool)Start(size int, inSize int, do func(interface{})interface{}, out chan interface{}) chan interface{} {
    waker := make(chan interface{}, inSize)
    self.objs = make([]*workUnit, size)
    wu := self.objs
    for i:=0;i<len(wu);i++{
        wu[i] = makeWorkUnit(i)
    }
    for _, obj := range wu {
        obj.Start(do, waker, out)
    }
    return waker
}

func (self *WorkUnitPool)Shutdown(){
    for _, obj := range self.objs {
        *obj.IsClosing <- true
        //fmt.Printf("Waiting for <%v>\n", obj.Name)
        obj.wg.Wait()
    }
    //? fmt.Println("Pool closed")
}


func (self *WorkUnitPool)PrintInfo() {
    for _, obj := range self.objs {
        fmt.Printf("<%v> is waken for %v time(s)\n", obj.Name, obj.RefCounter)
    }
}


const (
    DefaultWorkCount = 4
)
