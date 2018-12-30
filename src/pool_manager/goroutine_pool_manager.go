package pool_manager

import (
	"errors"
	"github.com/GanLuo96214/goroutine_pool/src/pool"
	"time"
)

type Interface interface {
	PoolManager() *pool.Status
}

var (
	pools map[string]*pool.Status
)

var addNameAlreadyBeUsed = errors.New("name already be used")

func Add(name string, p Interface) error {
	_, ok := pools[name]
	if ok {
		return addNameAlreadyBeUsed
	}

	pools[name] = p.PoolManager()

	return nil
}

// todo add release
//func Release(name string) {
//
//
//
//}

func SetExpectRunningCount(name string, count uint64) error {
	return pools[name].SetExpectRunningCount(count)
}
func GetExpectRunningCount(name string) uint64 {
	return pools[name].GetExpectRunningCount()
}

func SetDetectExpectDuration(name string, duration time.Duration) error {
	return pools[name].SetDetectExpectDuration(duration)
}
func GetDetectExpectDuration(name string) time.Duration {
	return pools[name].GetDetectExpectDuration()
}

func GetNowRunningCount(name string) uint64 {
	return pools[name].GetNowRunningCount()
}

func Info(name string) *pool.Status {
	return pools[name]
}

func All() map[string]*pool.Status {
	return pools
}
