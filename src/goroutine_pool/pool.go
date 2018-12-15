package goroutine_pool

import (
	"errors"
	"log"
	"math"
	"sync"
	"time"
)

type pool struct {
	expectRunningCount   uint64
	nowRunningCount      uint64
	nowRunningCountMutex sync.Mutex
	runFunc              func()
	detectExpectDuration time.Duration
	Index                uint64
	IndexMutex           sync.Mutex
	pool                 map[uint64]*bool
	poolMutex            sync.Mutex
	poolIndexChannel     chan uint64
}

var (
	poolSetFuncFuncIsNilError = errors.New("func is nil,the pool would not start")
)

func NewPool(
	expectRunningPoolCount uint64,
	runFunc func(),
) (gp *pool, err error) {

	gp = new(pool)

	// set default revise  running count
	err = gp.setDetectExpectDuration(time.Millisecond)
	if err != nil {
		return nil, err
	}

	// set expect running  count
	err = gp.setExpectRunningCount(expectRunningPoolCount)
	if err != nil {
		return nil, err
	}

	// init pool
	gp.pool = make(map[uint64]*bool, expectRunningPoolCount)

	if runFunc == nil {
		err = poolSetFuncFuncIsNilError
		return nil, err
	}

	gp.reviseRunningCount()

	return gp, err
}

var (
	poolSetExpectRunningCountMinCountError = errors.New("set count need >= 0") // count must >= 0
)

func (gp *pool) setExpectRunningCount(count uint64) (err error) {
	if count < 0 {
		err = poolSetExpectRunningCountMinCountError
		return err
	}

	gp.expectRunningCount = count

	return nil
}
func (gp *pool) getExpectRunningCount() uint64 {
	return gp.expectRunningCount
}

var (
	poolSetReviseCountMinDurationError = errors.New("min duration is millisecond") // min duration is millisecond because of cpu resource
)

func (gp *pool) setDetectExpectDuration(duration time.Duration) (err error) {
	if duration < time.Millisecond {
		err = poolSetReviseCountMinDurationError
		return
	}

	gp.detectExpectDuration = duration

	return
}
func (gp *pool) getDetectExpectDuration() time.Duration {
	return gp.detectExpectDuration
}

func (gp *pool) getNowRunningCount() uint64 {
	gp.nowRunningCountMutex.Lock()
	defer gp.nowRunningCountMutex.Unlock()
	return gp.nowRunningCount
}

func (gp *pool) incrNowRunningCount() {
	gp.nowRunningCountMutex.Lock()
	defer gp.nowRunningCountMutex.Unlock()
	gp.nowRunningCount++
}
func (gp *pool) decrNowRunningCount() {
	gp.nowRunningCountMutex.Lock()
	defer gp.nowRunningCountMutex.Unlock()
	gp.nowRunningCount--
}
func (gp *pool) newIndex() uint64 {
	gp.IndexMutex.Lock()
	defer gp.IndexMutex.Unlock()

	gp.Index++

	if gp.Index == math.MaxUint64 {
		gp.Index = 0
	}

	return gp.Index
}

func (gp *pool) reviseRunningCount() {
	go func() {
		defer gp.reviseRunningCount()

		if gp.expectRunningCount == gp.nowRunningCount {
			time.Sleep(gp.detectExpectDuration)
			return
		}

		if gp.expectRunningCount > gp.nowRunningCount {

			gp.addContainer()

		} else {

			gp.makeOneContainerToEnd()

		}

	}()
}

func (gp *pool) addContainer() {
	defer gp.poolMutex.Unlock()
	gp.poolMutex.Lock()

	index := gp.newIndex()
	isBreak := new(bool)

	gp.pool[index] = isBreak
	gp.poolIndexChannel <- index

	gp.container(isBreak, index)

}

func (gp *pool) makeOneContainerToEnd() {
	defer gp.poolMutex.Unlock()
	gp.poolMutex.Lock()

	index := <-gp.poolIndexChannel
	*gp.pool[index] = true
}

var (
	poolFunContainerIsBreakIsNilError = errors.New("isBreak is nil,this  would not start ")
)

func (gp *pool) container(isBreak *bool, index uint64) {
	defer gp.containerEnded(index)
	defer gp.decrNowRunningCount()

	if isBreak == nil {
		log.Println(poolFunContainerIsBreakIsNilError)
	}

	gp.incrNowRunningCount()

loop:
	gp.runFunc()

	if *isBreak {
		return
	}

	goto loop
}

func (gp *pool) containerEnded(index uint64) {
	defer gp.poolMutex.Unlock()
	gp.poolMutex.Lock()

	delete(gp.pool, index)

}
