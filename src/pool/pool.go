package pool

import (
	"errors"
	"sync"
	"time"
)

type pool struct {
	*Status

	reviseContainerRunningCountAsExpectCountMutex sync.Mutex

	runFunc func(containerIndex uint64)
}

var (
	newPoolRunFuncIsNil = errors.New("func is nil,the pool would not start")
)

func newPool(
	expectRunningCount uint64,
	runFunc func(containerIndex uint64),
) (p *pool, err error) {

	p = new(pool)
	p.Status = new(Status)

	// set default revise  running count
	err = p.SetDetectExpectDuration(defaultDetectExpectDuration)
	if err != nil {
		return nil, err
	}

	// set expect running  count
	err = p.SetExpectRunningCount(expectRunningCount)
	if err != nil {
		return nil, err
	}

	if runFunc == nil {
		err = newPoolRunFuncIsNil
		return nil, err
	}

	p.runFunc = runFunc

	// if GetNowRunningCount() < GetExpectRunningCount() then create containers
	go p.reviseContainerRunningCountAsExpectCount()

	return p, err
}

func (p *pool) containerStart(containerIndex uint64) {
	p.incrNowRunningCount()
	defer p.decrNowRunningCount()

	p.reviseContainerRunningCountAsExpectCountMutex.Unlock()

	p.runFunc(containerIndex)

	return
}

func (p *pool) reviseContainerRunningCountAsExpectCount() {
	for {
		p.reviseContainerRunningCountAsExpectCountMutex.Lock()

		if p.GetNowRunningCount() == p.GetExpectRunningCount() || p.GetNowRunningCount() > p.GetExpectRunningCount() {
			p.reviseContainerRunningCountAsExpectCountMutex.Unlock()
			time.Sleep(p.GetDetectExpectDuration())
			continue
		}

		go p.containerStart(p.newContainerIndex())
	}
}
