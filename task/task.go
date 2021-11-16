package task

import (
	"context"
	"errors"
	"stress-plan/helper"
	"stress-plan/logger"
	"stress-plan/output"
	"stress-plan/sender"
	"strings"
	"sync"
)

type Task struct {
	recWg  sync.WaitGroup
	sendWg sync.WaitGroup
	// 请求client
	sender sender.Interface
	// 请求结果格式化输出
	io     output.Interface
	logger logger.Logger
}

func NewTask(req sender.Interface, io output.Interface, logger logger.Logger) *Task {
	return &Task{
		sender: req,
		io:     io,
		logger: logger,
	}
}

// TODO 请求模版单独文件夹存放
const HttpMethodSplitFlag = "@"

// 请求模版，如果没有mock数据，那么和请求一致
type ReqTemplate struct {
	url     string
	headers map[string]string
	body    []byte
}

func NewReqTemplate(url string, headers map[string]string, body []byte) *ReqTemplate {
	return &ReqTemplate{
		url:     url,
		headers: headers,
		body:    body,
	}
}

// 由模版动态生成请求
func (r ReqTemplate) GenerateReq() (*sender.Request, error) {
	req := &sender.Request{
		Headers: make(map[string]string, len(r.headers)),
	}

	url := r.url
	isHttp := strings.Contains(url, "https://") || strings.Contains(url, "http://")
	if isHttp {
		req.Typ = "http"
		req.Method = "Get" // default method
		if i := strings.Index(url, HttpMethodSplitFlag); i != -1 {
			req.Method = url[:i]
			url = url[i+1:]
		}
	} else {
		return nil, errors.New("now not soupport this request type!")
	}

	// is wss

	req.Url = helper.PeplaceAllSpecialChars(url)
	bodyStr := helper.PeplaceAllSpecialChars(helper.UnsafeBytesToStr(r.body))
	req.Body = helper.UnsafeStrToBytes(bodyStr)
	for k, v := range r.headers {
		req.Headers[k] = helper.PeplaceAllSpecialChars(v)
	}

	return req, nil
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
	var gerr error

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
			gerr = t.sendManyRequests(context.Background(), sc, reqTemp, ch)
		}(sendCount)
	}

	if gerr != nil {
		t.logger.Log(logger.Error, "sendRequest", gerr.Error())
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
func (t *Task) sendManyRequests(ctx context.Context, sendCount int, reqTemp *ReqTemplate, result chan<- *sender.Result) error {
	for i := 0; i < sendCount; i++ {
		req, err := reqTemp.GenerateReq()
		if err != nil {
			return err
		}

		t.sender.Send(ctx, req, result)
	}
	return nil
}

// func GO(func() error) {
// 	defer
// }
