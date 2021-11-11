package sender

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

const statisticsDataDocs = `
=========== 统计结果文档 ==============
code 0 成功
`

type StatisticData struct {
	Concurrent int
	SuccessNum int
	FailureNum int
	TimeOutNum int
	// 总请求时间
	ReqTotalTime time.Duration
	// 实际请求消耗时间 = 总请求时间 / Concurrent
	ReqActualTime time.Duration
	// 平均请求时间
	AverageTime time.Duration
	// 最大、最小请求时间
	MaxTime time.Duration
	MinTime time.Duration
	QPS     float64
	Details string
}

// TODO 起一个协程每1s输出一下进度
// 输出统计过程，return 统计结果
func StatisticalResults(concurrent int, ch <-chan *Result) *StatisticData {
	var sd StatisticData
	var builder strings.Builder
	statusCodeMap := make(map[int]int, 3)
	statusCodes := make([]int, 0, 3)
	builder.Grow(1024)
	sd.MinTime = math.MaxInt32 * time.Millisecond

	sd.Concurrent = concurrent
	for result := range ch {
		sd.ReqTotalTime += result.UsedTime

		if sd.MaxTime < result.UsedTime {
			sd.MaxTime = result.UsedTime
		}
		if result.UsedTime < sd.MinTime {
			sd.MinTime = result.UsedTime
		}

		if result.StatusCode == 0 {
			sd.SuccessNum++
		} else {
			sd.FailureNum++
		}

		if result.IsTimeOut {
			sd.TimeOutNum++
		}

		if _, ok := statusCodeMap[result.StatusCode]; !ok {
			statusCodes = append(statusCodes, result.StatusCode)
		}
		statusCodeMap[result.StatusCode]++
	}
	if sd.SuccessNum+sd.FailureNum > 0 {
		sd.AverageTime = sd.ReqTotalTime / time.Duration(sd.SuccessNum+sd.FailureNum)
	}
	sd.ReqActualTime = sd.ReqTotalTime / time.Duration(sd.Concurrent)
	sd.QPS = (float64(sd.SuccessNum) + float64(sd.FailureNum)) / sd.ReqActualTime.Seconds()

	// 顺序的写入状态码信息
	sort.Ints(statusCodes)
	for _, code := range statusCodes {
		builder.WriteString(fmt.Sprintf("code: %d, count %d;\n", code, statusCodeMap[code]))
	}
	// builder.WriteString(statisticsDataDocs)
	sd.Details = builder.String()

	return &sd
}

// 生成统计过程中的信息
func StatisticsProcessInfo() string {
	return ""
}
