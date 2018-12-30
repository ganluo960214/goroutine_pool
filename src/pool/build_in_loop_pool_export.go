package pool

func NewBuildInLoopPool(
	expectRunningCount uint64,
	runFunc func(*bool, uint64),
) (p *buildInLoopPool, err error) {
	return newBuildInLoopPool(expectRunningCount, runFunc)
}
