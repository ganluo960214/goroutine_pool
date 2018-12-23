package pool

import (
	"errors"
	"sync"
	"time"
)

const (
	defaultDetectExpectDuration = time.Second
)

type pool struct {
	expectRunningCount   uint64
	nowRunningCount      uint64
	nowRunningCountMutex sync.Mutex
	detectExpectDuration time.Duration

	containerPrepareNext      chan *bool
	containerPrepareNextMutex sync.Mutex

	reviseContainerRunningCountAsExpectCountMutex sync.Mutex

	runFunc func()
}

var (
	newPoolRunFuncIsNil = errors.New("func is nil,the pool would not start")
)

func newPool(
	expectRunningCount uint64,
	runFunc func(),
) (p *pool, err error) {

	p = new(pool)

	// set default revise  running count
	err = p.setDetectExpectDuration(defaultDetectExpectDuration)
	if err != nil {
		return nil, err
	}

	// set expect running  count
	err = p.setExpectRunningCount(expectRunningCount)
	if err != nil {
		return nil, err
	}

	if runFunc == nil {
		err = newPoolRunFuncIsNil
		return nil, err
	}

	p.runFunc = runFunc

	p.containerPrepareNext = make(chan *bool)

	// if GetNowRunningCount() < GetExpectRunningCount() then release a container
	go p.reviseContainerRunningCountAsExpectCount()
	// if GetNowRunningCount() > GetExpectRunningCount() then release a container
	go p.reviseOverflowContainer()

	return p, err
}

var (
	setExpectRunningCountMinCountError = errors.New("set count need >= 0") // count must >= 0
)

func (p *pool) setExpectRunningCount(count uint64) (err error) {
	if count < 0 {
		err = setExpectRunningCountMinCountError
		return err
	}

	p.expectRunningCount = count

	return nil
}
func (p *pool) getExpectRunningCount() uint64 {
	return p.expectRunningCount
}

var (
	setDetectExpectDurationMinDurationError = errors.New("min duration is millisecond") // min duration is millisecond because of cpu resource
)

func (p *pool) setDetectExpectDuration(duration time.Duration) (err error) {
	if duration < time.Millisecond {
		err = setDetectExpectDurationMinDurationError
		return
	}

	p.detectExpectDuration = duration

	return
}
func (p *pool) getDetectExpectDuration() time.Duration {
	return p.detectExpectDuration
}

func (p *pool) incrNowRunningCount() {
	p.nowRunningCountMutex.Lock()
	defer p.nowRunningCountMutex.Unlock()

	p.nowRunningCount++
}
func (p *pool) decrNowRunningCount() {
	p.nowRunningCountMutex.Lock()
	defer p.nowRunningCountMutex.Unlock()

	p.nowRunningCount--
}
func (p *pool) getNowRunningCount() uint64 {
	p.nowRunningCountMutex.Lock()
	defer p.nowRunningCountMutex.Unlock()
	return p.nowRunningCount
}

func (p *pool) newContainerBreaker() *bool {
	return new(bool)
}

func (p *pool) containerStart(containerBreaker *bool) {
	p.incrNowRunningCount()
	defer p.decrNowRunningCount()

	p.reviseContainerRunningCountAsExpectCountMutex.Unlock()

	for *containerBreaker == false {
		p.runFunc()
		p.containerPrepareNextMutex.Lock()
		p.containerPrepareNext <- containerBreaker
	}

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

		go p.containerStart(p.newContainerBreaker())
	}
}

func (p *pool) reviseOverflowContainer() {
	for {

		containerBreaker := <-p.containerPrepareNext

		if p.GetNowRunningCount() > p.GetExpectRunningCount() {
			*containerBreaker = true
		}

		p.containerPrepareNextMutex.Unlock()

	}
}
