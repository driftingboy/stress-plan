package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"stress-plan/helper"
	"stress-plan/logger"
	"stress-plan/output"
	"stress-plan/output/std"
	"stress-plan/sender"
	"stress-plan/sender/http"
)

func main() {
	runtime.GOMAXPROCS(1)
	log := logger.NewStdLogger(os.Stdout, 4096)
	sender := http.NewSender(http.WithTTL(19*time.Second), http.WithLogger(log))
	out := &std.StdOut{}

	task := NewTask(true, sender, out, log)

	rt := &ReqTemplate{
		Method: "Get",
		Url:    "https://www.baidu.com/sugrec?prod=pc_his&from=pc_web&json=1&sid=1464_21098_31424_31341_31464_31229_30823_31163_31475&hisdata=&req=2&csor=0",
		Headers: map[string]string{
			"Accept":          "application/json, text/javascript, */*; q=0.01",
			"Connection":      "keep-alive",
			"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36",
			"Sec-Fetch-Site":  "same-origin",
			"Sec-Fetch-Mode":  "cors",
			"Sec-Fetch-Dest":  "empty",
			"Referer":         "https://www.baidu.com/",
			"Accept-Language": "zh-CN,zh;q=0.9",
		},
	}
	start := time.Now()
	fmt.Println("压测任务开始...")
	if err := task.Run(10, 400, rt); err != nil {
		log.Log(logger.Error, "统计错误", err.Error())
	}
	fmt.Printf("\n压测任务耗时：%v \n", time.Since(start))
	fmt.Println("压测任务结束!")
}

type Task struct {
	enableMock bool
	recWg      sync.WaitGroup
	sendWg     sync.WaitGroup
	// 请求client
	sender sender.Interface
	// 请求结果格式化输出
	io     output.Interface
	logger logger.Logger
}

func NewTask(enableMock bool, req sender.Interface, io output.Interface, logger logger.Logger) *Task {
	return &Task{
		enableMock: enableMock,
		sender:     req,
		io:         io,
		logger:     logger,
	}
}

// 请求模版，如果没有mock数据，那么和请求一致
type ReqTemplate struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    []byte
}

// 由模版动态生成请求
func (r ReqTemplate) GenerateReq() *sender.Request {
	req := &sender.Request{
		Headers: make(map[string]string, len(r.Headers)),
	}

	req.Method = r.Method
	req.Url = helper.PeplaceAllSpecialChars(r.Url)
	bodyStr := helper.PeplaceAllSpecialChars(helper.UnsafeBytesToStr(r.Body))
	req.Body = helper.UnsafeStrToBytes(bodyStr)
	for k, v := range r.Headers {
		req.Headers[k] = helper.PeplaceAllSpecialChars(v)
	}

	return req
}

// 直接转换为 sender.Request，不做任何处理
func (r ReqTemplate) Convert2Req() *sender.Request {
	return &sender.Request{
		Method:  r.Method,
		Url:     r.Url,
		Headers: r.Headers,
		Body:    r.Body,
	}
}

// 简单的分批任务
// TODO 重构为 协程池的模式，多个worker持续消费，保证并发度的稳定
func (t *Task) Run(concurrent, totalNum int, reqTemp *ReqTemplate) error {

	// 校验入参
	if concurrent < 1 || totalNum < 1 {
		return errors.New("concurrent and totalNum can not be less than 1!")
	}
	if reqTemp == nil {
		return errors.New("reqTemp can not be nil!")
	}

	ch := make(chan *sender.Result, 100)
	// 统计数据
	t.recWg.Add(1)
	var statisticsData *sender.StatisticData
	go func() {
		defer t.recWg.Done()
		statisticsData = sender.StatisticalResults(concurrent, ch)
	}()

	// 并发发送数据
	if concurrent > totalNum {
		concurrent = totalNum
	}
	perNum := totalNum / concurrent
	remainder := totalNum % concurrent
	for i := 0; i < concurrent; i++ {
		t.sendWg.Add(1)
		sendCount := perNum
		if i == concurrent-1 {
			sendCount += remainder
		}
		go func(sc int) {
			defer t.sendWg.Done()
			t.sendManyRequests(context.Background(), sc, reqTemp, ch)
		}(sendCount)
	}

	// 等待发送结束
	t.sendWg.Wait()
	// 关闭通道，通知Statistical结束循环
	close(ch)
	// 等待统计结束，输出统计结果
	t.recWg.Wait()
	return t.io.Write(statisticsData)
}

// TODO 传入请求数据，每次发送前通过 mock_parser对数据处理
func (t *Task) sendManyRequests(ctx context.Context, sendCount int, reqTemp *ReqTemplate, result chan<- *sender.Result) {
	for i := 0; i < sendCount; i++ {
		var req *sender.Request
		if t.enableMock {
			req = reqTemp.GenerateReq()
		} else {
			req = reqTemp.Convert2Req()
		}
		t.sender.Send(ctx, req, result)
	}
}

// func GO(func() error) {
// 	defer
// }
