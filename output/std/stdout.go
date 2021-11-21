package std

import (
	"errors"
	"fmt"
	"os"
	"stress-plan/helper"
	"stress-plan/sender"
)

type StdOut struct {
}

func (s *StdOut) Write(data *sender.StatisticData) (err error) {
	if data == nil {
		return errors.New("统计数据为空！")
	}

	fmt.Fprintf(os.Stdout, `
	Stressing test in %d goroutine %d request

		Success: %d failed: %d timeout: %d

		Req total Duration：%v s
		Req Actual Time：%v s
		avg    max    min
		%v ms  %v ms  %v ms

		RequestsPeerSec: %v
		TransferPeerSec: %v

	%d requests completed, %s Flow transfer in %v seconds
	Detail：%s`,
		data.Concurrent, data.SuccessNum+data.FailureNum,
		data.SuccessNum, data.FailureNum, data.TimeOutNum,
		data.ReqTotalTime.Seconds(), data.ReqActualTime.Seconds(), data.AverageTime.Milliseconds(), data.MaxTime.Milliseconds(), data.MinTime.Milliseconds(),
		data.RequestPeerSec, helper.BytesAddUnit(int64(data.TransferBytesPeerSec)),
		data.SuccessNum+data.FailureNum, helper.BytesAddUnit(data.TransferBytes), data.ReqActualTime.Seconds(),
		data.Details)
	return nil
}
