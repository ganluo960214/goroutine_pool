package pool_manager

import "github.com/GanLuo96214/goroutine_pool/src/pool"

type PoolManager interface {
	PoolManager() *pool.Status
}
