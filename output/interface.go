package output

import (
	"io"
	"stress-plan/sender"
)

type Interface interface {
	// ReaderWithUnmarshal
	// WriterWithMarshal
	Write(data *sender.StatisticData) (err error)
}

type StrandTestDatas struct{}

type ReaderWithUnmarshal interface {
	io.Reader
	Unmarshal([]byte) (*StrandTestDatas, error)
}

type WriterWithMarshal interface {
	io.Writer
	Marshal(*StrandTestDatas) ([]byte, error)
}
