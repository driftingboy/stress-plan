package sender

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStatisticalResults(t *testing.T) {
	resultSet := []*Result{
		{Gid: 1, IsTimeOut: false, UsedTime: 200 * time.Millisecond, ResponseBytes: 2000, StatusCode: 0},
		{Gid: 1, IsTimeOut: false, UsedTime: 1000 * time.Millisecond, ResponseBytes: 200, StatusCode: 401},
		{Gid: 2, IsTimeOut: true, UsedTime: 1200 * time.Millisecond, ResponseBytes: 100, StatusCode: 500},
		{Gid: 2, IsTimeOut: false, UsedTime: 300 * time.Millisecond, ResponseBytes: 100, StatusCode: 0},
		{Gid: 3, IsTimeOut: false, UsedTime: 100 * time.Millisecond, ResponseBytes: 100, StatusCode: 0},
	}

	ch := make(chan *Result, 100)
	go func(ch chan *Result) {
		defer close(ch)
		for _, r := range resultSet {
			ch <- r
		}
	}(ch)
	got := StatisticalResults(2, ch)
	want := &StatisticData{
		Concurrent:    2,
		SuccessNum:    3,
		FailureNum:    2,
		TimeOutNum:    1,
		ReqTotalTime:  2800 * time.Millisecond,
		ReqActualTime: 1400 * time.Millisecond,
		AverageTime:   560 * time.Millisecond,
		MaxTime:       1200 * time.Millisecond,
		MinTime:       100 * time.Millisecond,
		QPS:           5 / 1.4,
		Details:       "code: 0, count 3;\ncode: 401, count 1;\ncode: 500, count 1;\n",
	}
	assert.Equal(t, want, got)

}
