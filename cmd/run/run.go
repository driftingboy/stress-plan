package run

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"stress-plan/logger"
	"stress-plan/output/std"
	"stress-plan/sender/http"
	"stress-plan/task"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// 这里并发数可能存在歧义，两种理解
// 一：起了多少线程、协程，就是多少并发，并发为10
// 二：根据实际情况，比如10个协程，一条请求平均250ms完成，那么一个协程1s处理4条请求，10个协程单位时间1s内能处理40个请求，tps为40
// 如果把1s当作单位时间，并发就是指在`同一段时间内`同时处理,也有人理解 并发为40

// 用户数和并发数的关系，一般并发数 = 同时在人数 * 1% ～ 5%（当然只是针对大部分场景）
// 也就是你想测试 1w 用户同时在线，-c 为100左右即可

var (
	cpu         int8
	concurrency uint64
	totalNumber uint64

	// TODO 一些没有生效 链接设置
	conn    int
	isHttp2 bool
	// TODO keepalive 需要在 sender 中每次请求判断，如果是则加上 req.Close = true
	isLongConn bool
	timeout    uint64 // s

	// 请求设置
	url     string
	body    string // TODO 支持传文件
	headers map[string]string

	// 统计设置
	// TODO
	follow bool
	// TODO
	code = 200 // 成功状态码, 默认200
)

func RunCMD() *cobra.Command {
	runCmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"execute"},
		Short:   "execute stressing test task",
		Long:    "execute stressing test task",
		RunE: func(_ *cobra.Command, args []string) error {
			fmt.Printf("run args %+v", args)
			return run()
		},
	}

	flags := runCmd.Flags()

	flags.Int8VarP(&cpu, "cpu", "", 1, "压测客户端使用逻辑cpu数, cpu < 0 则都使用")
	flags.Uint64VarP(&concurrency, "concurrency", "c", 1, "并发数")
	flags.Uint64VarP(&totalNumber, "number", "n", 1, "请求数总量")

	flags.StringVarP(&url, "url", "u", "", "压测地址, e.g. 'GET@https://www.baidu.com/'")
	flags.StringVarP(&body, "body", "b", "", "请求体，用于http请求")
	flags.StringToStringVarP(&headers, "header", "H", make(map[string]string), "自定义头信息 e.g. -H 'Content-Type: application/json'")

	flags.BoolVarP(&isHttp2, "http2", "", false, "是否开http2.0")
	flags.BoolVarP(&isLongConn, "keepalive", "k", true, "是否开启长连接")
	flags.IntVarP(&conn, "conn", "", 5, "单个host最大连接数")
	flags.Uint64VarP(&timeout, "timeout", "", 30, "请求超时时间")

	flags.BoolVarP(&follow, "follow", "f", false, "是否动态打印统计信息")

	return runCmd
}

func run() error {
	if err := verifyFlag(); err != nil {
		return err
	}

	// new task
	usablecpu := runtime.NumCPU()
	cpuInt := int(cpu)
	if cpuInt > usablecpu {
		cpuInt = usablecpu
	}
	runtime.GOMAXPROCS(cpuInt)

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
		return err
	}

	fmt.Printf("\n压测任务耗时：%v \n", time.Since(start))
	fmt.Println("压测任务结束!")
	return nil
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
