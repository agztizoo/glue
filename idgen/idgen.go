package idgen

import (
	"strconv"
	"time"
)

// IDGenerator 代表 ID 生成器.
type IDGenerator interface {

	// GenID 生成字符串 ID.
	GenID() (string, error)

	// MustGenID 生成字符串 ID 或 panic.
	MustGenID() string

	// GenIntID 生成整型 ID.
	GenIntID() (int64, error)

	// MustGenIntID 生成整型 ID 或 panic.
	MustGenIntID() int64

	// TimeOfID 提取 ID 包含的时间.
	TimeOfID(string) (time.Time, error)

	// TimeOfIntID 提取 ID 包含的时间.
	TimeOfIntID(int64) (time.Time, error)

	// MustTimeOfID 提取 ID 包含的时间或 panic.
	MustTimeOfID(string) time.Time

	// MustTimeOfIntID 提取 ID 包含的时间或 panic.
	MustTimeOfIntID(int64) time.Time
}

// New 创建 IDGenerator.
func New(
	// ID 生成函数.
	idf func() (int64, error),
	// ID 转换 时间函数.
	tmf func(id int64) (time.Time, error),
) IDGenerator {
	return &idg{idf: idf, tmf: tmf}
}

type idg struct {
	// ID 生成函数.
	idf func() (int64, error)
	// ID 转换 时间函数.
	tmf func(id int64) (time.Time, error)
}

// GenID 生成字符串 ID.
func (i *idg) GenID() (string, error) {
	id, err := i.idf()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

// GenIntID 生成整型 ID.
func (i *idg) GenIntID() (int64, error) {
	return i.idf()
}

// MustGenID 生成字符串 ID 或 panic.
func (i *idg) MustGenID() string {
	id, err := i.GenID()
	if err != nil {
		panic(err)
	}
	return id
}

// MustGenIntID 生成整型 ID 或 panic.
func (i *idg) MustGenIntID() int64 {
	id, err := i.GenIntID()
	if err != nil {
		panic(err)
	}
	return id
}

// TimeOfIntID 提取 ID 包含的时间.
func (i *idg) TimeOfIntID(id int64) (time.Time, error) {
	return i.tmf(id)
}

// TimeOfID 提取 ID 包含的时间.
func (i *idg) TimeOfID(s string) (time.Time, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return i.TimeOfIntID(id)
}

// MustTimeOfID 提取 ID 包含的时间或 panic.
func (i *idg) MustTimeOfID(s string) time.Time {
	t, err := i.TimeOfID(s)
	if err != nil {
		panic(err)
	}
	return t
}

// MustTimeOfIntID 提取 ID 包含的时间或 panic.
func (i *idg) MustTimeOfIntID(id int64) time.Time {
	t, err := i.TimeOfIntID(id)
	if err != nil {
		panic(err)
	}
	return t
}
