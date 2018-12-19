package pool

import (
	"container/list"
	"errors"
	"log"
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

	containerStatusList      *list.List
	containerStatusListMutex sync.Mutex

	reviseContainerRunningCountMutex sync.Mutex

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

	p.containerStatusList = list.New()

	// start revise container running count
	go p.reviseContainerRunningCount()

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

func (p *pool) newContainerStatus() *list.Element {
	isBreak := new(bool)

	containerStatus := p.pushNewContainerStatusInto(isBreak)

	return containerStatus
}

func (p *pool) pushNewContainerStatusInto(isBreak *bool) *list.Element {
	defer p.containerStatusListMutex.Unlock()
	p.containerStatusListMutex.Lock()
	containerStatus := p.containerStatusList.PushBack(isBreak)

	return containerStatus
}

var (
	containerContainerStatusTypeError = errors.New("p.containerStatusList type error(require *bool)")
)

func (p *pool) containerStart(containerStatus *list.Element) {
	defer p.containerEnd(containerStatus)
	defer p.decrNowRunningCount()
	defer p.reviseContainerRunningCountMutex.Unlock()

	var isBreak *bool
	switch v := containerStatus.Value.(type) {
	case *bool:
		isBreak = v
	default:
		err := containerContainerStatusTypeError
		log.Println(err)
		return
	}

	p.incrNowRunningCount()
	p.reviseContainerRunningCountMutex.Unlock()

	for *isBreak == false {
		p.runFunc()
	}

	return
}
func (p *pool) containerEnd(containerStatus *list.Element) {
	defer p.containerStatusListMutex.Unlock()
	p.containerStatusListMutex.Lock()
	p.containerStatusList.Remove(containerStatus)
}

func (p *pool) reviseContainerRunningCount() {
	for {

		p.reviseContainerRunningCountMutex.Lock() // this lock will release in p.container method

		if p.GetExpectRunningCount() == p.GetNowRunningCount() {
			p.reviseContainerRunningCountMutex.Unlock()
			time.Sleep(p.GetDetectExpectDuration())
			continue
		}

		if p.GetNowRunningCount() < p.GetExpectRunningCount() {
			containerStatus := p.newContainerStatus()
			go p.containerStart(containerStatus)
		} else { // release a container

			var isBreak *bool
			switch v := p.containerStatusList.Front().Value.(type) {
			case *bool:
				isBreak = v
			default:
				log.Fatal("p.containerStatusList find no *bool in list.List.Element. and I don't know why. Pool already will try stop all container. shutdown now.")
				return
			}
			*isBreak = true

		}

	}
}
