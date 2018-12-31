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

	p, err := pool.NewPool(10, func(containerIndex uint64) {
		fmt.Printf("#%d routinue executed\n", containerIndex)
		time.Sleep(time.Second)
		// end of execution
		// return <-  same as end of execution
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
	err = p.SetDetectExpectDuration(100)
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
