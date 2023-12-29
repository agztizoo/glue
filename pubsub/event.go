package pubsub

import "time"

var (
	// MessageTimeFunc 生成消息时间.
	MessageTimeFunc = time.Now
)

const (
	// MessageHeaderAggregateType 预定义 Header Key.
	// 领域消息聚合根类型.
	MessageHeaderAggregateType = "glue:pubsub:aggregate_type"

	// MessageHeaderAggregateID 领域消息聚合根 ID.
	MessageHeaderAggregateID = "glue:pubsub:aggregate_id"

	// MessageHeaderTenantID 触发消息用户所在租户 ID.
	MessageHeaderTenantID = "glue:pubsub:tenant_id"

	// MessageHeaderUserID 触发消息用户 ID.
	MessageHeaderUserID = "glue:pubsub:user_id"

	// 兼容：请使用  MessageHeader_xxx.

	// EventHeaderAggregateType 领域消息聚合根类型.
	EventHeaderAggregateType = MessageHeaderAggregateType

	// EventHeaderAggregateID 领域消息聚合根 ID.
	EventHeaderAggregateID = MessageHeaderAggregateID

	// EventHeaderTenantID 触发消息用户所在租户 ID.
	EventHeaderTenantID = MessageHeaderTenantID

	// EventHeaderUserID 触发消息用户 ID.
	EventHeaderUserID = MessageHeaderUserID
)

// NewMessage 创建消息.
func NewMessage(id string, payload interface{}, createAt time.Time) *Message {
	return &Event{
		id:       id,
		headers:  make(map[string]interface{}),
		payload:  payload,
		createAt: createAt,
	}
}

// Message 定义 Publisher 转发给 Subscriber 的消息.
type Message struct {
	// 消息 ID.
	id string

	// 消息头信息.
	headers map[string]interface{}

	// Publisher 中发布的消息内容.
	payload interface{}

	// 消息创建时间.
	createAt time.Time
}

// Event 兼容.
type Event = Message

// NewEvent 兼容.
var NewEvent = NewMessage

// GetID 返回消息 ID.
func (msg *Message) GetID() string {
	return msg.id
}

// GetCreateTime 返回消息创建时间.
func (msg *Message) GetCreateTime() time.Time {
	return msg.createAt
}

// GetAggregateType 获取领域消息所属聚合类型.
func (msg *Message) GetAggregateType() string {
	v := msg.headers[EventHeaderAggregateType]
	if v == nil {
		return ""
	}

	typ, ok := v.(string)
	if !ok {
		return ""
	}
	return typ
}

func (msg *Message) GetAggregateID() string {
	v := msg.headers[EventHeaderAggregateID]
	if v == nil {
		return ""
	}

	id, ok := v.(string)
	if !ok {
		return ""
	}
	return id
}

// GetTenantID 获取触发消息用户所属租户 ID.
func (msg *Message) GetTenantID() string {
	v := msg.headers[EventHeaderTenantID]
	if v == nil {
		return ""
	}

	id, ok := v.(string)
	if !ok {
		return ""
	}
	return id
}

// GetUserID 获取触发消息用户 ID.
func (msg *Message) GetUserID() string {
	v := msg.headers[EventHeaderUserID]
	if v == nil {
		return ""
	}

	id, ok := v.(string)
	if !ok {
		return ""
	}
	return id
}

// SetHeader 设置消息头信息.
func (msg *Message) SetHeader(key string, value interface{}) {
	msg.headers[key] = value
}

// GetHeader 获取消息头信息.
func (msg *Message) GetHeader(key string) interface{} {
	return msg.headers[key]
}

// GetHeaders 获取消息所有头信息.
//
// 返回当前头信息的浅拷贝.
func (msg *Message) GetHeaders() map[string]interface{} {
	hdrs := make(map[string]interface{})
	for k, v := range msg.headers {
		hdrs[k] = v
	}
	return hdrs
}

// SetPayload 设置消息内容.
func (msg *Message) SetPayload(payload interface{}) {
	msg.payload = payload
}

// GetPayload 获取消息内容.
func (msg *Message) GetPayload() interface{} {
	return msg.payload
}
