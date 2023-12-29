package pubsub

import "context"

type DomainEvent interface {
	EventType() string
}

type EventPublisher interface {
	// Publish 实现领域事件发布.
	Publish(ctx context.Context, aggregateType, aggregateID string, events ...DomainEvent)
}
