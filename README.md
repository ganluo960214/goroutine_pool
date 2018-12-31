# goroutine_pool

goroutine_pool is a framework for management goroutines.
 
## features 
- all kind of pool can manipulate goroutines running count as u expected.
    - [NewPool](#newpool)
    - [NewBuildInLoopPool](#newbuildinlooppool)
- all kind of pool has stateful container(container: pool's function container).
    - [NewPool](#newpool)
    - [NewBuildInLoopPool](#newbuildinlooppool)
- a small pool manager
    - [PoolManager](#poolmanager)


## Contents
- [NewPool](#newpool)
- [NewBuildInLoopPool](#newbuildinlooppool)
- [PoolManager](#poolmanager)

## NewPool

- quick start

```go
// main.go
package main

import (
    "fmt"
    "github.com/GanLuo96214/goroutine_pool/src/pool"
    "log"
    "sync"
    "time"
)

func main() {

    wg := sync.WaitGroup{}

    var err error

    p,err := pool.NewPool(10, func(containerIndex uint64) {
        fmt.Printf("#%d routinue executed\n",containerIndex)
        time.Sleep(time.Second)
        // end of execution
        // return // <-  same as end of execution
    })
    if err != nil {
        log.Fatal(err)
    }

    // set pool expect running count
    err = p.SetExpectRunningCount(100)
    if err != nil {
        log.Fatal(err)
    }

    // set pool detect duration
    err = p.SetDetectExpectDuration(100)
    if err != nil {
        log.Fatal(err)
    }

    // get pool expect running count
    expectRunningCount := p.GetExpectRunningCount()
    fmt.Printf("%d expect running count\n",expectRunningCount)
    // get pool detect duration(min detect is time.Millisecond)
    detectDuration := p.GetDetectExpectDuration()
    fmt.Printf("%f second detect",float64(detectDuration / time.Second))
    // get pool now running count
    nowRunningCount := p.GetNowRunningCount()
    fmt.Printf("%d now running count\n",nowRunningCount)

    wg.Add(1)
    wg.Wait()
    
    /**
    output:
    #1 routinue executed
    #2 routinue executed
    #3 routinue executed
    #4 routinue executed
    #5 routinue executed
    #6 routinue executed
    #7 routinue executed
    #8 routinue executed
    #9 routinue executed
    #10 routinue executed
    ...after 1s
    #11 routinue executed
    #12 routinue executed
    #13 routinue executed
    #14 routinue executed
    #16 routinue executed
    #15 routinue executed
    #18 routinue executed
    #17 routinue executed
    #19 routinue executed
    #20 routinue executed
    ...after 1s
    #21 routinue executed
    #22 routinue executed
    #23 routinue executed
    #24 routinue executed
    #25 routinue executed
    #26 routinue executed
    #27 routinue executed
    #28 routinue executed
    #29 routinue executed
    #30 routinue executed
    ...
     */
    
}
```

- notices
  - when a pool's function end of execution will decrement 1 running count then will start a new one container by `DetectExpectDuration` will increment 1 running count
  - when a pool's function end of execution the `containerIndex` will be gone with it, the new one container  will got a new `containerIndex`(1 to math.MaxUint64, when arrived math.MaxUint64 next will be 1).
  - u also can write a endless loop in function but recommend use [NewBuildInLoopPool](#newbuildinlooppool)
  - also can ignored status use it as stateless.

## NewBuildInLoopPool

- quick start

```go
package main

import (
	"fmt"
	"github.com/GanLuo96214/goroutine_pool/src/pool"
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {

    wg := sync.WaitGroup{}
    
    var err error
    
    p, err := pool.NewBuildInLoopPool(10, func(containerEnd func(), containerIndex uint64) {
        // this function running in endless loop
        fmt.Printf("#%d routinue executed\n", containerIndex)
        time.Sleep(time.Second)
    
    
        r := rand.Intn(3)
        if r == 3 { // only 3 will end loop.
            containerEnd() // <- call this function end loop to start a new container with new containerIndex
            return // immediately end of execution also u can running belong code
        }
    
        // end of execution
        // return <-  end of execution
        // but not will end the loop
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // set pool expect running count
    err = p.SetExpectRunningCount(10)
    //err = p.SetExpectRunningCount(0)
    //err = p.SetExpectRunningCount(1000)
    if err != nil {
        log.Fatal(err)
    }
    
    // set pool detect duration(min detect is time.Millisecond)
    err = p.SetDetectExpectDuration(time.Second)
    if err != nil {
        log.Fatal(err)
    }
    
    // get pool expect running count
    expectRunningCount := p.GetExpectRunningCount()
    fmt.Printf("%d expect running count\n", expectRunningCount)
    // get pool detect duration(min detect is time.Millisecond)
    detectDuration := p.GetDetectExpectDuration()
    fmt.Printf("%f second detect", float64(detectDuration/1000000000))
    // get pool now running count
    nowRunningCount := p.GetNowRunningCount()
    fmt.Printf("%d now running count\n", nowRunningCount)
    
    wg.Add(1)
    wg.Wait()
    
    /**
    output:
    #1 routinue executed
    #2 routinue executed
    #3 routinue executed
    #4 routinue executed
    #5 routinue executed
    #6 routinue executed
    #7 routinue executed
    #8 routinue executed
    #9 routinue executed
    #10 routinue executed
    ...after 1s
    #1 routinue executed
    #2 routinue executed
    #3 routinue executed
    #4 routinue executed
    #5 routinue executed
    #6 routinue executed
    #7 routinue executed
    #8 routinue executed
    #9 routinue executed
    #10 routinue executed
    ...after 1s
    #1 routinue executed
    #2 routinue executed
    #3 routinue executed
    #4 routinue executed
    #5 routinue executed
    #6 routinue executed
    #7 routinue executed
    #8 routinue executed
    #9 routinue executed
    #10 routinue executed
    ...        
    */
}
```

- notices
  - when a pool's function end of execution will running again with same state
  - when a pool's function inside called containerEnd() and end of execution will decrement 1 running count then will start a new one container by `DetectExpectDuration` will increment 1 running count
  - when a pool's function inside called containerEnd() and end of execution the `containerIndex` will be gone with it, the new one container  will got a new `containerIndex`(1 to math.MaxUint64, when arrived math.MaxUint64 next will be 1).
  - also can ignored status use it as stateless.

## PoolManager(todo)
