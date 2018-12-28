package pool

import "time"

func NewBuildInLoopPool(
	expectRunningCount uint64,
	runFunc func(*bool, uint64),
) (p *buildInLoopPool, err error) {
	return newBuildInLoopPool(expectRunningCount, runFunc)
}

func (p *buildInLoopPool) SetExpectRunningCount(count uint64) (err error) {
	return p.setExpectRunningCount(count)
}
func (p *buildInLoopPool) GetExpectRunningCount() uint64 {
	return p.expectRunningCount
}

func (p *buildInLoopPool) SetDetectExpectDuration(duration time.Duration) (err error) {
	return p.setDetectExpectDuration(duration)
}
func (p *buildInLoopPool) GetDetectExpectDuration() time.Duration {
	return p.getDetectExpectDuration()
}

func (p *buildInLoopPool) GetNowRunningCount() uint64 {
	return p.getNowRunningCount()
}
