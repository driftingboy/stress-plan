/**
 * @description: 简单替换特殊字符
 * TODO 范围, 如 @Int[0,10]
 * TODO 类型，如 @UnixNano @RFC3339
 */
package helper

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	guid "github.com/google/uuid"
)

const SpecFlag = "@"

type Int string

func (i Int) Mock() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return strconv.Itoa(r.Intn(math.MaxInt32))
}

func (i Int) Name() string {
	return "@Int"
}

func (i Int) Len() int {
	// int32的最大长度10 2147483647
	return 10
}

type UUID string

func (u UUID) Mock() string {
	return guid.New().String()
}

func (u UUID) Name() string {
	return "@UUID"
}

func (u UUID) Len() int {
	// uuid string xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	return 36
}

type Unix string

func (u Unix) Mock() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func (u Unix) Name() string {
	// unix时间戳的最大长度
	return "@Unix"
}

func (u Unix) Len() int {
	return 10
}

type Mocker interface {
	// 生成mock数据
	Mock() string
	Name() string
	Len() int
}

var (
	integer = new(Int)
	unix    = new(Unix)
	uuid    = new(UUID)

	intName, intNameLen   = integer.Name(), len(intName)
	unixName, unixNameLen = unix.Name(), len(unixName)
	uuidName, uuidNameLen = uuid.Name(), len(uuidName)
)

// 解析特殊字符 special characters
func PeplaceAllSpecialChars(data string) string {
	if strings.IndexByte(data, SpecFlag[0]) < 0 {
		return data
	}
	newLen, replaceCount := CalcNewStringLength(data)

	// 拼接新字符串
	newString := make([]byte, newLen)
	w := 0
	start := 0
	for i := 0; i < replaceCount; i++ {
		index, mocker := IndexSpecialChars(data[start:])
		if index < 0 {
			break
		}
		w += copy(newString[w:], data[start:start+index])
		w += copy(newString[w:], mocker.Mock()[:])
		start += index + len(mocker.Name())
	}
	w += copy(newString[w:], data[start:])

	// todo unsafe convert
	return string(newString[0:w])
}

// 寻找有没有特殊替换，有则返回下标，无则返回-1
func IndexSpecialChars(data string) (index int, m Mocker) {
	dataLen := len(data)
	for i := 0; i < dataLen; i++ {
		hasSpecFlag := (data[i] == '@')
		if hasSpecFlag {
			if isInt := (i+intNameLen <= dataLen && data[i:i+intNameLen] == intName); isInt {
				return i, integer
			}
			if isUnix := (i+unixNameLen <= dataLen && data[i:i+unixNameLen] == unixName); isUnix {
				return i, unix
			}
			if isUUID := (i+uuidNameLen <= dataLen && data[i:i+uuidNameLen] == uuidName); isUUID {
				return i, uuid
			}
		}
	}
	return -1, nil
}

func CalcNewStringLength(old string) (newLen int, replaceCount int) {
	// 计算增加或减少了多少长度，以及出现多少次替换
	addLen := 0
	replaceCount = 0
	for i := 0; i < len(old); i++ {
		index, mocker := IndexSpecialChars(old[i:])
		if index < 0 {
			break
		}
		nameLen := len(mocker.Name())
		i += index + nameLen
		addLen += (mocker.Len() - nameLen)
		replaceCount++
	}

	return len(old) + addLen, replaceCount
}
