package http

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"time"

	"stress-plan/helper"
	"stress-plan/logger"
	"stress-plan/sender"
)

// go generate: impl -dir="./request/http"  "s *Sender" request.Interface
type Sender struct {
	// 在程序里控制超时
	// http本身的超时只能控制从tcp建立链接到接收完response的时间，发送request的时间没法控制
	ttl    time.Duration
	cli    *http.Client
	logger logger.Logger
}

// 默认为一个timeout=30s的长链接
func NewSender(opts ...Option) *Sender {
	s := &Sender{
		ttl:    30 * time.Second,
		cli:    DefaultLongConnClient,
		logger: logger.NewStdLogger(os.Stdout, 4096),
	}

	for _, o := range opts {
		o(s)
	}
	return s
}

// 发送请求，result用于接收请求结果
func (s *Sender) Send(ctx context.Context, req *sender.Request, result chan<- *sender.Result) {
	// 处理超时情况
	ctx, cancel := context.WithTimeout(ctx, s.ttl)
	defer cancel()

	// 构建请求
	body := bytes.NewReader(req.Body)
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.Url, body)
	if err != nil {
		// TODO 需要退出通道，监听退出。gorutione 同时退出
		s.logger.Log(logger.Error, "err", err.Error())
		return
	}
	req.Headers["Content-Type"] = "application/x-www-form-urlencoded; charset=utf-8"
	for k, v := range req.Headers {
		req.Headers[k] = v
	}

	// 结果返回
	var (
		isTimeOut bool
		usedTime  time.Duration
		resBytes  int64
		code      int = -1 // 表示发生失败
	)
	startTime := time.Now()
	resp, err := s.cli.Do(httpReq)
	usedTime = time.Since(startTime)

	if err != nil {
		s.logger.Log(logger.Error, "clientErr", err.Error())
		if e, ok := err.(interface{ Timeout() bool }); ok && e.Timeout() {
			isTimeOut = true
		} else {
			s.logger.Log(logger.Error, "clientErr", err.Error())
		}
	}
	if resp != nil {
		resBytes = resp.ContentLength
		if resp.StatusCode == 200 {
			code = 0
		} else {
			code = resp.StatusCode
		}
	}

	res := &sender.Result{
		Gid:           helper.GetGoidSlowly(),
		UsedTime:      usedTime,
		ResponseBytes: resBytes,
		StatusCode:    code,
		IsTimeOut:     isTimeOut,
	}
	result <- res
}
