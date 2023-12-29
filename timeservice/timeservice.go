package timeservice

import "time"

// TimeService 定义时间服务.
type TimeService interface {
	// Now 返回当前时间.
	Now() time.Time
}

// NewTimeService 创建基于机器时间的时间服务.
func NewTimeService() TimeService {
	return &timeService{}
}

type timeService struct {
}

// Now 返回当前时间.
func (_ *timeService) Now() time.Time {
	return time.Now()
}
