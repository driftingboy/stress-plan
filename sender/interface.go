package sender

import (
	"context"
	"errors"
	"time"

	"stress-plan/helper"
)

type Interface interface {
	// 发送请求，result用于接收请求结果
	Send(context context.Context, req *Request, result chan<- *Result)
	// IsVaildate() bool
	// Format(string) string
}

type Request struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    []byte
}

func NewRequest(method, url string, headers map[string]string, body string) *Request {
	return &Request{
		Method:  method,
		Url:     url,
		Headers: headers,
		Body:    helper.UnsafeStrToBytes(body),
	}
}

func (r *Request) Vaildate() error {
	if r == nil {
		return errors.New("request can not be nil!")
	}
	return nil
}

type Result struct {
	// goroutine id
	Gid int64
	// 消耗时长
	UsedTime time.Duration
	// 结果数据大小
	ResponseBytes int64
	// 是否超时
	IsTimeOut bool
	// 状态码
	StatusCode int
	// 格外的信息描述
	// ExtraMessage string
}
