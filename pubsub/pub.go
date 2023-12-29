package pubsub

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/glue/idgen"
	"github.com/glue/transaction"
)

// Publisher 定义“领域事件”发布器.
//
// 调度顺序：
//   1.同步调度在异步调度前.
//   2.逐个 Subscriber 对 events 依次调度.
//     for _, sub := range subscribes {
//         for _, event := events {
//             sub.OnEvent(ctx, event)
//         }
//     }
// 错误处理:
//   1.同步调度(Subscriber 需实现 SynchronousSubscriber):
//   ⚠️  注意: 事务处理. 除非明确知道注意事项, 否则不要实现。
//     a.Subscriber 之间的错误互相影响, 出错则后续调度终止。
//     b.当 Subscriber 实现了 EventGroupSubscriber 则 event group 整体调度.
//     c.当 event group 中出现 event 调度失败, 则 event group 后续事件调度终止.
//   2.异步调度(默认):
//     a.Subscriber 之间的错误互不影响.
//     b.当 Subscriber 实现了 EventGroupSubscriber 则 event group 整体调度.
//     c.当 event group 中出现 event 调度失败, 则 event group 后续事件调度终止.
//
// 事务处理:
//   1. 当 Subscriber 是同步调度，则继承 ctx 中的事务上下文.
//   2. 当 Subscriber 是异步调度, 则逃脱 ctx 中的事务进行调度.
type Publisher interface {
	// Publish 发布事件.
	//
	// 不支持指定聚合类型和聚合 ID.
	Publish(ctx context.Context, event ...interface{}) error
}

// MakeFromContextRender 创建通过 context 渲染事件的渲染器.
//
// 从 context 提取 key-value 并渲染到 Event Header, key 为空字符串时忽略.
func MakeFromContextRender(kvFuns ...func(context.Context) (key string, val interface{})) func(context.Context, *Message) *Message {
	return func(ctx context.Context, msg *Message) *Message {
		for _, fn := range kvFuns {
			key, val := fn(ctx)
			if key == "" {
				continue
			}
			msg.SetHeader(key, val)
		}
		return msg
	}
}

// NewPublisher 创建“领域事件”发布器.
//
// render 对发布前的事件进行渲染，如: 填充 TenantID、UserID.
// render 为 nil 时不进行渲染.
func NewPublisher(idg idgen.IDGenerator, tm transaction.Manager, render func(context.Context, *Message) *Message, subs ...Subscriber) Publisher {
	pub := &publisher{idg: idg, tm: tm, render: render}
	for _, sub := range subs {
		if _, ok := sub.(SynchronousSubscriber); ok {
			pub.syncs = append(pub.syncs, sub)
		} else {
			pub.asyncs = append(pub.asyncs, sub)
		}
	}
	return pub
}

type publisher struct {
	idg    idgen.IDGenerator
	tm     transaction.Manager
	render func(context.Context, *Message) *Message

	syncs  []Subscriber
	asyncs []Subscriber
}

// Publish 实现领域事件发布.
func (p *publisher) Publish(ctx context.Context, event ...interface{}) error {
	msgs, err := p.toMessages(ctx, event)
	if err != nil {
		return err
	}
	return p.publish(ctx, msgs...)
}

func (p *publisher) publish(ctx context.Context, msg ...*Message) error {
	// 同步事件调度.
	for _, sub := range p.syncs {
		if err := p.scheduleSync(ctx, sub, p.filterSupportedMessages(sub, msg)); err != nil {
			return err
		}
	}
	// 注册事务成功回调.
	registered := p.tm.OnCommitted(ctx, func(ctx context.Context) {
		// 异步事件调度 - 事务 commit 成功后.
		for _, sub := range p.asyncs {
			p.scheduleAsync(ctx, sub, p.filterSupportedMessages(sub, msg))
		}
	})
	// 事务回调注册失败.
	if !registered {
		// 异步事件调度 - 不在事务内,立即触发.
		for _, sub := range p.asyncs {
			p.scheduleAsync(ctx, sub, p.filterSupportedMessages(sub, msg))
		}
	}
	return nil
}

func (p *publisher) toMessages(ctx context.Context, events []interface{}) ([]*Message, error) {
	msgs := make([]*Event, 0, len(events))
	for _, event := range events {
		ev, err := p.toMessage(ctx, event)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, ev)
	}
	return msgs, nil
}

func (p *publisher) toMessage(ctx context.Context, event interface{}) (*Event, error) {
	msg, ok := event.(*Message)
	if !ok {
		id, err := p.idg.GenID()
		if err != nil {
			return nil, err
		}
		msg = NewEvent(id, event, MessageTimeFunc())
	}
	if p.render != nil {
		msg = p.render(ctx, msg)
	}
	return msg, nil
}

func (p *publisher) filterSupportedMessages(sub Subscriber, msgs []*Message) []*Message {
	var ms []*Message
	for _, msg := range msgs {
		if !sub.IsSupported(msg) {
			continue
		}
		ms = append(ms, msg)
	}
	return ms
}

// scheduleSync 同步调度事件.
func (p *publisher) scheduleSync(ctx context.Context, sub Subscriber, msgs []*Message) error {
	if len(msgs) <= 0 {
		return nil
	}
	// 支持 event group 调度.
	if s, ok := sub.(MessageGroupSubscriber); ok {
		return s.OnEventGroup(ctx, msgs...)
	}
	for _, event := range msgs {
		if err := sub.OnEvent(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// scheduleAsync 异步调度事件.
func (p *publisher) scheduleAsync(ctx context.Context, sub Subscriber, msgs []*Message) {
	if len(msgs) <= 0 {
		return
	}
	p.tm.EscapeTransaction(ctx, func(ctx context.Context) error {
		go p.scheduleIgnoreError(ctx, sub, msgs)
		return nil
	})
}

// scheduleIgnoreError 顺序调度 event group, 出错则终止 event group 中后续事件调度.
func (p *publisher) scheduleIgnoreError(ctx context.Context, sub Subscriber, msgs []*Message) {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithContext(ctx).Errorf("[glue][publisher] schedule subscriber: %s panic: %v", sub.Name(), err)
		}
	}()

	// 支持 event group 调度.
	if s, ok := sub.(MessageGroupSubscriber); ok {
		if err := s.OnEventGroup(ctx, msgs...); err != nil {
			logrus.WithContext(ctx).Errorf("[glue][publisher] failed to schedule event group subscriber: %s error: %v", sub.Name(), err)
		}
		return
	}

	for _, msg := range msgs {
		if err := sub.OnEvent(ctx, msg); err != nil {
			logrus.WithContext(ctx).Errorf("[glue][publisher] failed to schedule subscriber: %s event: %s error: %v", sub.Name(), msg.GetID(), err)
		}
	}
}

// NewEventPublisher 创建“领域事件”发布器.
//
// render 对发布前的事件进行渲染，如: 填充 TenantID、UserID.
// render 为 nil 时不进行渲染.
func NewEventPublisher(idg idgen.IDGenerator, tm transaction.Manager, render func(context.Context, *Message) *Message, subs ...Subscriber) EventPublisher {
	pub := &eventPublisher{idg: idg, tm: tm, render: render}
	for _, sub := range subs {
		if _, ok := sub.(SynchronousSubscriber); ok {
			pub.syncs = append(pub.syncs, sub)
		} else {
			pub.asyncs = append(pub.asyncs, sub)
		}
	}
	return pub
}

type eventPublisher publisher

func (ep *eventPublisher) Publish(ctx context.Context, aggregateType, aggregateID string, events ...DomainEvent) {
	p := (*publisher)(ep)
	evs, err := p.toMessages(ctx, ep.domainEventToInterface(events))
	if err != nil {
		panic(err)
	}
	for _, ev := range evs {
		ev.SetHeader(EventHeaderAggregateType, aggregateType)
		ev.SetHeader(EventHeaderAggregateID, aggregateID)
	}
	if err := p.publish(ctx, evs...); err != nil {
		panic(err)
	}
}

func (ep *eventPublisher) domainEventToInterface(events []DomainEvent) []interface{} {
	evs := make([]interface{}, 0, len(events))
	for _, ev := range events {
		evs = append(evs, ev)
	}
	return evs
}
