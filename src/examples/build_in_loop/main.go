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
			return         // immediately end of execution also u can running belong code
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
