package pool_manager

import "github.com/GanLuo96214/goroutine_pool/src/pool"

func init() {
	pools = make(map[string]*pools.Status)
}
