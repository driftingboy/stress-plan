package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"stress-plan/logger"
	"stress-plan/output/std"
	"stress-plan/sender/http"
	"stress-plan/task"
)

// flag
var (
	cpu         int
	concurrency uint64
	totalNumber uint64

	// 链接设置
	conn       int
	isHttp2    bool
	isLongConn bool
	timeout    uint

	// 请求设置
	url     string
	body    string
	headers mapper = make(mapper)

	// 统计设置
	code = 200 // 成功状态码, 默认200
)

var _ flag.Value = (*mapper)(nil)

type mapper map[string]string

func (m mapper) String() string {
	return fmt.Sprint(map[string]string(m))
}

func (m mapper) Set(s string) error {
	kv := strings.SplitN(s, ":", 2)
	if len(kv) < 2 {
		return errors.New("header format err, such as `key:value`")
	}
	m[kv[0]] = kv[1]
	return nil
}

func init() {
	flag.IntVar(&cpu, "cpu", 1, "cpu核数，< 1则都使用")
	flag.Uint64Var(&concurrency, "c", 1, "并发数")
	flag.Uint64Var(&totalNumber, "n", 1, "请求数总量")

	flag.StringVar(&url, "u", "", "压测地址, 目前支持http/https; 如果是http/https,需要加上请求方法, e.g. 'Get|https://www.baidu.com/'")
	flag.Var(&headers, "H", "自定义头信息 如:-H 'Content-Type: application/json'")
	flag.StringVar(&body, "b", "", "HTTP POST方式传送数据")

	flag.IntVar(&conn, "conn", 1, "单个host最大连接数")
	flag.BoolVar(&isHttp2, "http2", false, "是否开http2.0")
	flag.BoolVar(&isLongConn, "keepalive", true, "是否开启长连接")
	flag.UintVar(&timeout, "timeout", 30, "请求超时时间（秒）")

	flag.IntVar(&code, "code", 200, "请求成功的状态码")

	// 解析参数
	flag.Parse()
}

func verifyFlag() error {
	if concurrency < 1 {
		return errors.New("concurrency can not be less than 1!")
	}
	if totalNumber < 1 {
		return errors.New("totalNumber can not be less than 1!")
	}
	if conn < 1 {
		return errors.New("max conn can not be less than 1!")
	}
	if timeout < 1 {
		return errors.New("timeout can not be less than 1!")
	}
	if strings.TrimSpace(url) == "" {
		return errors.New("url can not be less than empty!  use -u")
	}

	return nil
}

func main() {
	if err := verifyFlag(); err != nil {
		fmt.Println(err.Error())
		return
	}

	// new task
	usablecpu := runtime.NumCPU()
	if cpu > usablecpu {
		cpu = usablecpu
	}
	runtime.GOMAXPROCS(cpu)

	log := logger.WithPrefix(logger.NewStdLogger(os.Stdout, 4096), logger.DefaultCaller, logger.DefaultTimer)
	opts := []http.Option{
		http.WithLogger(log),
		http.WithTTL(time.Duration(timeout) * time.Second),
	}
	if isHttp2 {
		opts = append(opts, http.WithHttp2())
	}
	if !isLongConn {
		opts = append(opts, http.WithDisableKeepAlive())
	}

	sender := http.NewSender(opts...)
	out := &std.StdOut{}

	t := task.NewTask(sender, out, log)

	// task start
	fmt.Println("压测任务开始...")

	start := time.Now()
	rt := task.NewReqTemplate(url, headers, []byte(body))
	if err := t.Run(int(concurrency), int(totalNumber), rt); err != nil {
		log.Log(logger.Error, "统计错误", err.Error())
	}

	fmt.Printf("\n压测任务耗时：%v \n", time.Since(start))
	fmt.Println("压测任务结束!")
}
