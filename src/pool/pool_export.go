package pool

import "time"

func NewPool(
	expectRunningCount uint64,
	runFunc func(),
) (p *pool, err error) {
	return newPool(expectRunningCount, runFunc)
}

func (p *pool) SetExpectRunningCount(count uint64) (err error) {
	return p.setExpectRunningCount(count)
}
func (p *pool) GetExpectRunningCount() uint64 {
	return p.expectRunningCount
}

func (p *pool) SetDetectExpectDuration(duration time.Duration) (err error) {
	return p.setDetectExpectDuration(duration)
}
func (p *pool) GetDetectExpectDuration() time.Duration {
	return p.getDetectExpectDuration()
}

func (p *pool) GetNowRunningCount() uint64 {
	return p.getNowRunningCount()
}
