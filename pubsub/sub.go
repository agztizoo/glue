package pubsub

import "context"

// Subscriber 定义“领域事件”订阅器.
//
// 参考 Publisher 说明.
//
// Subscriber 默认为异步订阅, 异步订阅的事件消费在发布方的事务成功提交后.
//
// 如果要实现为同步订阅, 比如做 Event-Sourcing. 需要 Subscriber 实现 SynchronousSubscriber.
// ⚠️  注意: 由于事件在事务内发布，同步订阅需要确保没有事务外操作。如：RPC.
type Subscriber interface {
	// Name 返回 Subscriber 名称.
	Name() string

	// IsSupported 判断订阅器是否支持该领域事件.
	IsSupported(msg *Message) bool

	// OnEvent 处理订阅事件.
	//
	// IsAsync 为 true 时在新 goroutine 进行调度.
	//
	// Publish 中的 event group 被逐个调度, group 中一个事件调度失败则 group 中后
	// 续事件调度终止.
	//
	// 如果需要 event group 被整体调度需要实现 EventGroupSubscriber.
	OnEvent(ctx context.Context, msg *Message) error
}

// SynchronousSubscriber 表示 Subscriber 是同步订阅器.
//
// ⚠️  注意: 事务处理. 除非明确知道注意事项, 否则不要实现。
type SynchronousSubscriber interface {
	// SynchronousSubscriber 空实现, 仅用于同步订阅器标记.
	SynchronousSubscriber()
}

// MessageGroupSubscriber 表示 Subscriber 是否支持“领域事件”组整体订阅.
//
// 当 Subscriber 同时实现了 MessageGroupSubscriber 则支持 event group 整体调度.
type MessageGroupSubscriber interface {
	// OnEventGroup 处理订阅事件.
	//
	// Publish 中的 event group 被整体调度.
	OnEventGroup(ctx context.Context, msg ...*Message) error
}

// NewSubscriber 创建事件订阅器.
//
// 默认为异步事件订阅，在新的 goroutine 消费事件.
func NewSubscriber(name string, supportFn func(*Message) bool, onEventFn func(context.Context, *Message) error) Subscriber {
	return &subscriber{
		name:      name,
		supportFn: supportFn,
		onEventFn: onEventFn,
	}
}

type subscriber struct {
	name      string
	supportFn func(*Message) bool
	onEventFn func(context.Context, *Message) error
}

// Name 返回 Subscriber 名称.
func (s *subscriber) Name() string {
	return s.name
}

// IsSupported 判断订阅器是否支持该领域事件.
func (s *subscriber) IsSupported(msg *Message) bool {
	return s.supportFn(msg)
}

// OnEvent 处理订阅事件.
func (s *subscriber) OnEvent(ctx context.Context, msg *Message) error {
	return s.onEventFn(ctx, msg)
}
