package pool

import (
	"errors"
	"math"
	"sync"
	"time"
)

const (
	defaultDetectExpectDuration = time.Second
)

type Status struct {
	expectRunningCount   uint64
	nowRunningCount      uint64
	nowRunningCountMutex sync.Mutex
	containerIndex       uint64
	containerIndexMutex  sync.Mutex
	detectExpectDuration time.Duration
}

var (
	setExpectRunningCountMinCountError = errors.New("set count need >= 0") // count must >= 0
)

func (s *Status) SetExpectRunningCount(count uint64) (err error) {
	if count < 0 {
		err = setExpectRunningCountMinCountError
		return err
	}

	s.expectRunningCount = count

	return nil
}
func (s *Status) GetExpectRunningCount() uint64 {
	return s.expectRunningCount
}

var (
	setDetectExpectDurationMinDurationError = errors.New("min duration is millisecond") // min duration is millisecond because of cpu resource
)

func (s *Status) SetDetectExpectDuration(duration time.Duration) (err error) {
	if duration < time.Millisecond {
		err = setDetectExpectDurationMinDurationError
		return
	}

	s.detectExpectDuration = duration

	return
}
func (s *Status) GetDetectExpectDuration() time.Duration {
	return s.detectExpectDuration
}

func (s *Status) incrNowRunningCount() {
	s.nowRunningCountMutex.Lock()
	defer s.nowRunningCountMutex.Unlock()

	s.nowRunningCount++
}
func (s *Status) decrNowRunningCount() {
	s.nowRunningCountMutex.Lock()
	defer s.nowRunningCountMutex.Unlock()

	s.nowRunningCount--
}
func (s *Status) GetNowRunningCount() uint64 {
	s.nowRunningCountMutex.Lock()
	defer s.nowRunningCountMutex.Unlock()
	return s.nowRunningCount
}

func (s *Status) newContainerIndex() uint64 {
	s.containerIndexMutex.Lock()
	defer s.containerIndexMutex.Unlock()

	if s.containerIndex == math.MaxUint64 {
		s.containerIndex = 0
	}

	s.containerIndex++

	return s.containerIndex
}

func (s *Status) newContainerBreaker() *bool {
	return new(bool)
}

func (s *Status) PoolManager() *Status {

	if s == nil {
		return nil
	}

	return s
}
