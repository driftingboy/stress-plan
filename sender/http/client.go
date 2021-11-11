package http

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

// 定制 http client

var (
	DefaultShortConnClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			// DisableCompression: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
			MaxIdleConns:        0,
			MaxConnsPerHost:     100,
			MaxIdleConnsPerHost: 100,
		},
	}
	DefaultLongConnClient = &http.Client{
		// 设置超时时间（包括从连接(Dial)到读完response body）
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				// 限制建立连接(Dial)的时间
				Timeout: 30 * time.Second,
				// 设置keepalive失效时间
				KeepAlive: 60 * time.Second,
			}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
			MaxIdleConns:        0, // 0 表示不限制
			MaxConnsPerHost:     100,
			MaxIdleConnsPerHost: 100,
		},
	}
	DefaultLongConnHttp2 = GetLongConnHttp2()
)

func GetLongConnHttp2() *http.Client {
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
		MaxIdleConns:        0,
		MaxConnsPerHost:     100,
		MaxIdleConnsPerHost: 100,
	}
	_ = http2.ConfigureTransport(tr)
	return &http.Client{Timeout: 60 * time.Second, Transport: tr}
}
