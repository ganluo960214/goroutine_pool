package goroutine_pool

import "time"

// export start
func (gp *pool) SetExpectRunningCount(count uint64) (err error) {
	return gp.setExpectRunningCount(count)
}
func (gp *pool) GetExpectRunningCount() uint64 {
	return gp.expectRunningCount
}

func (gp *pool) SetDetectExpectDuration(duration time.Duration) (err error) {
	return gp.setDetectExpectDuration(duration)
}
func (gp *pool) GetDetectExpectDuration() time.Duration {
	return gp.detectExpectDuration
}

func (gp *pool) GetNowRunningCount() uint64 {
	return gp.getNowRunningCount()
}

// export end
