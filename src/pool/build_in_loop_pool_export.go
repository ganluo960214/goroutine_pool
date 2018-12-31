package pool

func NewBuildInLoopPool(
	expectRunningCount uint64,
	runFunc func(containerEnd func(), containerIndex uint64),
) (p *buildInLoopPool, err error) {
	return newBuildInLoopPool(expectRunningCount, runFunc)
}
