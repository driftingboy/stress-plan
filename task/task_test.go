package task

import (
	"fmt"
	"os"
	"runtime"
	"stress-plan/logger"
	"stress-plan/output/std"
	"stress-plan/sender/http"
	"testing"
	"time"
)

func Test_Task(t *testing.T) {
	runtime.GOMAXPROCS(1)
	log := logger.WithPrefix(logger.NewStdLogger(os.Stdout, 4096), logger.DefaultCaller, logger.DefaultTimer)
	// TODO 2.flag 使用包优化，stp -h --help, stp run -c 20 -d 100s -u xxxxx,
	sender := http.NewSender(http.WithTTL(18*time.Second), http.WithLogger(log))
	out := &std.StdOut{}

	task := NewTask(sender, out, log)

	rt := &ReqTemplate{
		url:     "Get@https://stackoverflow.com/questions/43321894/context-timeout-implementation-on-every-request-using-golang",
		headers: map[string]string{
			// "Accept":          "application/json, text/javascript, */*; q=0.01",
			// "Connection":      "keep-alive",
			// "User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36",
			// "Sec-Fetch-Site":  "same-origin",
			// "Sec-Fetch-Mode":  "cors",
			// "Sec-Fetch-Dest":  "empty",
			// "Referer":         "https://www.baidu.com/",
			// "Accept-Language": "zh-CN,zh;q=0.9",
		},
	}
	start := time.Now()
	fmt.Println("压测任务开始...")
	if err := task.Run(1, 4, rt); err != nil {
		log.Log(logger.Error, "统计错误", err.Error())
	}
	fmt.Printf("\n压测任务耗时：%v \n", time.Since(start))
	fmt.Println("压测任务结束!")
}
