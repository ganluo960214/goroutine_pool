package pool

func NewPool(
	expectRunningCount uint64,
	runFunc func(containerIndex uint64),
) (p *pool, err error) {
	return newPool(expectRunningCount, runFunc)
}
