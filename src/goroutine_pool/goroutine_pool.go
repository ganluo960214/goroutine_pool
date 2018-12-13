package goroutine_pool

import (
	"errors"
	"log"
	"math"
	"sync"
	"time"
)

type GoroutinePool struct {
	expectRunningGoroutineCount uint64
	nowRunningGoroutineCount uint64
	nowRunningGoroutineCountMutex sync.Mutex
	runFunc func()
	detectExpectDuration time.Duration
	goroutineIndex uint64
	goroutineIndexMutex sync.Mutex
	pool map[uint64]*bool
	poolMutex sync.Mutex
	poolGoroutineIndexChannel chan uint64
}

var (
	GoroutinePoolSetFuncFuncIsNilError = errors.New("func is nil,the pool would not start")
)
func (g *GoroutinePool) SetFuncAndPoolCount(
	poolCount uint64,
	runFunc func(),
) (err error) {

	// set default revise goroutine running count
	err = g.SetReviseGoroutineCountDuration(time.Millisecond)
	if err != nil {
		return err
	}

	// set expect running goroutine count
	err = g.SetExpectRunningGoroutineCount(poolCount)
	if err != nil {
		return err
	}

	// init pool
	g.pool = make(map[uint64]*bool,poolCount)

	if runFunc == nil {
		err = GoroutinePoolSetFuncFuncIsNilError
		return err
	}

	return err
}


var (
	GoroutinePoolSetExpectRunningGoroutineCountMinCountError = errors.New("set count need >= 0") // count must >= 0
)
func (g *GoroutinePool) SetExpectRunningGoroutineCount(count uint64) (err error) {
	if count < 0 {
		err = GoroutinePoolSetExpectRunningGoroutineCountMinCountError
		return err
	}

	g.expectRunningGoroutineCount = count

	return nil
}
func (g *GoroutinePool) GetExpectRunningGoroutineCount() uint64 {
	return g.expectRunningGoroutineCount
}


var (
	GoroutinePoolSetReviseGoroutineCountMinDurationError = errors.New("min duration is millisecond") // min duration is millisecond because of cpu resource
)
func (g *GoroutinePool) SetReviseGoroutineCountDuration(duration time.Duration) (err error) {
	if duration < time.Millisecond {
		err = GoroutinePoolSetReviseGoroutineCountMinDurationError
		return
	}

	g.detectExpectDuration = duration

	return
}
func (g *GoroutinePool) GetDetectExpectDuration() time.Duration {
	return g.detectExpectDuration
}

func (g *GoroutinePool) incrNowRunningGoroutineCount()  {
	g.nowRunningGoroutineCountMutex.Lock()
	defer g.nowRunningGoroutineCountMutex.Unlock()
	g.nowRunningGoroutineCount++
}
func (g *GoroutinePool) decrNowRunningGoroutineCount()  {
	g.nowRunningGoroutineCountMutex.Lock()
	defer g.nowRunningGoroutineCountMutex.Unlock()
	g.nowRunningGoroutineCount--
}
func (g *GoroutinePool) newGoroutineIndex() uint64 {
	g.goroutineIndexMutex.Lock()
	defer g.goroutineIndexMutex.Unlock()

	g.goroutineIndex++

	if g.goroutineIndex == math.MaxUint64 {
		g.goroutineIndex = 0
	}

	return g.goroutineIndex
}

func (g *GoroutinePool) reviseGoroutineRunningCount() {
	go func() {

		if g.expectRunningGoroutineCount == g.nowRunningGoroutineCount {
			time.Sleep(g.detectExpectDuration)
			go g.reviseGoroutineRunningCount()
			return
		}

		if g.expectRunningGoroutineCount > g.nowRunningGoroutineCount {

			g.addContainer()

		} else {

			g.setAContainerToEnd()

		}

		go g.reviseGoroutineRunningCount()
		return
	}()
}

func (g *GoroutinePool) addContainer()  {
	defer g.poolMutex.Unlock()
	g.poolMutex.Lock()

	index := g.newGoroutineIndex()
	isBreak := new(bool)

	g.pool[index] = isBreak
	g.poolGoroutineIndexChannel<-index

	g.container(isBreak,index)

}

func (g *GoroutinePool) setAContainerToEnd()  {
	defer g.poolMutex.Unlock()
	g.poolMutex.Lock()

	index := <-g.poolGoroutineIndexChannel
	*g.pool[index] = true
}


var (
	GoroutinePoolFunContainerIsBreakIsNilError = errors.New("isBreak is nil,this goroutine would not start ")
)
func (g *GoroutinePool) container(isBreak *bool,index uint64) {
	defer g.containerEnded(index)
	defer g.decrNowRunningGoroutineCount()

	if isBreak == nil {
		log.Println(GoroutinePoolFunContainerIsBreakIsNilError)
	}

	g.incrNowRunningGoroutineCount()

	loop:
		g.runFunc()

	if *isBreak {
		return
	}

	goto loop
}


func (g *GoroutinePool) containerEnded(index uint64)  {
	defer g.poolMutex.Unlock()
	g.poolMutex.Lock()

	delete(g.pool,index)

}
