package std

import (
	"errors"
	"fmt"
	"os"
	"stress-plan/sender"
)

type StdOut struct {
}

func (s *StdOut) Write(data *sender.StatisticData) (err error) {
	if data == nil {
		return errors.New("统计数据为空！")
	}
	fmt.Fprintf(os.Stdout, `
	==================== 统计结果 =================
	并发数：%d

	成功：%d
	失败：%d
	超时：%d

	请求总耗时：%v s
	请求实际时长：%v s
	平均时长：%v ms
	最大时长：%v ms
	最小时长：%v ms

	qps：%v
	详情：%s`,
		data.Concurrent,
		data.SuccessNum, data.FailureNum, data.TimeOutNum,
		data.ReqTotalTime.Seconds(), data.ReqActualTime.Seconds(), data.AverageTime.Milliseconds(),
		data.MaxTime.Milliseconds(), data.MinTime.Milliseconds(),
		data.QPS, data.Details)
	return nil
}
