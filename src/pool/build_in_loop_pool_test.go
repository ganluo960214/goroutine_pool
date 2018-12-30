package pool

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestNewBuildInLoopPool(t *testing.T) {

	var err error

	_, err = NewBuildInLoopPool(
		1,
		func(containerBreaker *bool, containerIndex uint64) {
			time.Sleep(time.Second)
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewBuildInLoopPool(
		1,
		nil,
	)
	if err != newBuildInLoopPoolRunFuncIsNil {
		t.Fatal(err)
	}

}

func TestBuildInLoopPool_SetDetectExpectDuration(t *testing.T) {
	var err error
	p, err := NewBuildInLoopPool(
		1,
		func(containerBreaker *bool, containerIndex uint64) {
			time.Sleep(time.Second)
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	var detectExpectDuration time.Duration

	// second
	detectExpectDuration = time.Second
	err = p.SetDetectExpectDuration(detectExpectDuration)
	if err != nil {
		t.Fatal(err)
	}

	// active detect expect duration min duration error
	detectExpectDuration = time.Microsecond
	err = p.SetDetectExpectDuration(detectExpectDuration)
	if err != setDetectExpectDurationMinDurationError {
		t.Fatal(err)
	}

	// hours
	detectExpectDuration = time.Hour
	err = p.SetDetectExpectDuration(detectExpectDuration)
	if err != nil {
		t.Fatal(err)
	}

}

var (
	TestBuildInLoopPoolGetDetectExpectDurationDefaultValueError           = errors.New("detect expect duration default value error(should be default value)")
	TestBuildInLoopPoolGetDetectExpectDurationShouldBePreviousSetDuration = errors.New("detect expect duration value should be previous set duration")
)

func TestBuildInLoopPool_GetDetectExpectDuration(t *testing.T) {

	var err error
	p, err := NewBuildInLoopPool(
		1,
		func(containerBreaker *bool, containerIndex uint64) {
			time.Sleep(time.Second)
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if p.GetDetectExpectDuration() != defaultDetectExpectDuration {
		t.Fatal(TestBuildInLoopPoolGetDetectExpectDurationDefaultValueError)
	}

	var detectExpectDuration time.Duration

	// second
	detectExpectDuration = time.Second
	err = p.SetDetectExpectDuration(detectExpectDuration)
	if err != nil {
		t.Fatal(err)
	}

	// active detect expect duration min duration error
	detectExpectDuration = time.Microsecond
	err = p.SetDetectExpectDuration(detectExpectDuration)
	if err != setDetectExpectDurationMinDurationError {
		t.Fatal(err)
	}

	// GetDetectExpectDuration should be previous set(because previous set got error)
	if p.GetDetectExpectDuration() != time.Second {
		t.Fatal(TestBuildInLoopPoolGetDetectExpectDurationShouldBePreviousSetDuration)
	}

}

func TestBuildInLoopPool_SetExpectRunningCount(t *testing.T) {
	var (
		err                error
		expectRunningCount uint64 = 1
	)

	p, err := NewBuildInLoopPool(
		expectRunningCount,
		func(containerBreaker *bool, containerIndex uint64) {
			time.Sleep(time.Second)
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	expectRunningCount = 0
	err = p.SetExpectRunningCount(expectRunningCount)
	if err != nil {
		t.Fatal(err)
	}

}

var (
	TestBuildInLoopPoolGetExpectRunningCountNotMatchSetCountError = errors.New("expect running count not match the set count")
)

func TestBuildInLoopPool_GetExpectRunningCount(t *testing.T) {
	var (
		err                error
		expectRunningCount uint64 = 1
	)

	p, err := NewBuildInLoopPool(
		expectRunningCount,
		func(containerBreaker *bool, containerIndex uint64) {
			time.Sleep(time.Second)
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	expectRunningCount = 0
	err = p.SetExpectRunningCount(expectRunningCount)
	if err != nil {
		t.Fatal(err)
	}
	if p.GetExpectRunningCount() != expectRunningCount {
		t.Fatal(TestBuildInLoopPoolGetExpectRunningCountNotMatchSetCountError)
	}

}

var (
	TestBuildInLoopPoolGetNowRunningCountExpectRunningCountNotEqualSetRunningCount = errors.New("expect running count not equal set expect running count")
	TestBuildInLoopPoolGetNowRunningCountNotEqualExpectRunningCount                = errors.New("running count not equal expect running count")
)

func TestBuildInLoopPool_GetNowRunningCount(t *testing.T) {
	var (
		i uint64 = 0
	)

	{
		var (
			expectRunningCount uint64 = 1
			wg                        = sync.WaitGroup{}
		)
		for i = 0; i < expectRunningCount; i++ {
			wg.Add(1)
		}

		p, err := NewBuildInLoopPool(
			expectRunningCount,
			func(containerBreaker *bool, containerIndex uint64) {

				wg.Done()
				time.Sleep(time.Hour)

			},
		)
		if err != nil {
			t.Fatal(err)
		}

		wg.Wait()

		if expectRunningCount != p.GetExpectRunningCount() {
			t.Fatal(TestBuildInLoopPoolGetNowRunningCountExpectRunningCountNotEqualSetRunningCount)
		}
		if expectRunningCount != p.GetNowRunningCount() {
			t.Fatal(TestBuildInLoopPoolGetNowRunningCountNotEqualExpectRunningCount)
		}

	}

	{
		var (
			expectRunningCount uint64 = 100
			wg                        = sync.WaitGroup{}
		)

		for i = 0; i < expectRunningCount; i++ {
			wg.Add(1)
		}

		p, err := NewBuildInLoopPool(
			expectRunningCount,
			func(containerBreaker *bool, containerIndex uint64) {

				wg.Done()
				time.Sleep(time.Hour)

			},
		)
		if err != nil {
			t.Fatal(err)
		}

		wg.Wait()

		if expectRunningCount != p.GetExpectRunningCount() {
			t.Fatal(TestBuildInLoopPoolGetNowRunningCountExpectRunningCountNotEqualSetRunningCount)
		}
		if expectRunningCount != p.GetNowRunningCount() {
			t.Fatal(TestBuildInLoopPoolGetNowRunningCountNotEqualExpectRunningCount)
		}

		//
		expectRunningCount = 150
		for i = 0; i < expectRunningCount-p.GetNowRunningCount(); i++ {
			wg.Add(1)
		}

		err = p.SetExpectRunningCount(expectRunningCount)
		if err != nil {
			t.Fatal(err)
		}
		wg.Wait()

		if expectRunningCount != p.GetExpectRunningCount() {
			t.Fatal(TestBuildInLoopPoolGetNowRunningCountExpectRunningCountNotEqualSetRunningCount)
		}
		if expectRunningCount != p.GetNowRunningCount() {
			t.Fatal(TestBuildInLoopPoolGetNowRunningCountNotEqualExpectRunningCount)
		}

	}

	{
		var (
			expectRunningCount uint64 = 100
		)

		p, err := NewBuildInLoopPool(
			expectRunningCount,
			func(containerBreaker *bool, containerIndex uint64) {

				time.Sleep(time.Millisecond)

			},
		)
		if err != nil {
			t.Fatal(err)
		}

		//
		if expectRunningCount != p.GetExpectRunningCount() {
			t.Fatal(TestBuildInLoopPoolGetNowRunningCountExpectRunningCountNotEqualSetRunningCount)
		}
		for p.GetNowRunningCount() != p.GetExpectRunningCount() {
			time.Sleep(time.Millisecond)
		}

		//
		expectRunningCount = 150
		err = p.SetExpectRunningCount(expectRunningCount)
		if err != nil {
			t.Fatal(err)
		}
		if expectRunningCount != p.GetExpectRunningCount() {
			t.Fatal(TestBuildInLoopPoolGetNowRunningCountExpectRunningCountNotEqualSetRunningCount)
		}
		for p.GetNowRunningCount() != p.GetExpectRunningCount() {
			time.Sleep(time.Millisecond)
		}

		//
		expectRunningCount = 70
		err = p.SetExpectRunningCount(expectRunningCount)
		if err != nil {
			t.Fatal(err)
		}
		if expectRunningCount != p.GetExpectRunningCount() {
			t.Fatal(TestBuildInLoopPoolGetNowRunningCountExpectRunningCountNotEqualSetRunningCount)
		}
		for p.GetNowRunningCount() != p.GetExpectRunningCount() {
			time.Sleep(time.Millisecond)
		}

	}

}
