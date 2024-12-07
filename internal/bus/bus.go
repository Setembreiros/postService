package bus

import (
	"context"
	"encoding/json"
)

//go:generate mockgen -source=bus.go -destination=mock/bus.go

type Event struct {
	Type string
	Data []byte
}

type EventBus struct {
	subscribers map[string][]chan<- Event
	externalBus ExternalBus
}

type EventSubscription struct {
	EventType string
	Handler   EventHandler
}

type ExternalBus interface {
	Publish(event *Event) error
}

type EventHandler interface {
	Handle(event []byte)
}

func NewEventBus(externalBus ExternalBus) *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan<- Event),
		externalBus: externalBus,
	}
}

func (eb *EventBus) PublishLocal(event Event) {
	subscriberChannels := eb.subscribers[event.Type]

	for _, subscriberChannel := range subscriberChannels {
		subscriberChannel <- event
	}
}

func (eb *EventBus) Subscribe(subscription *EventSubscription, ctx context.Context) {
	subscriptionChan := make(chan Event)
	eb.subscribers[subscription.EventType] = append(eb.subscribers[subscription.EventType], subscriptionChan)
	go subscription.handle(subscriptionChan, ctx)
}

func (eb *EventBus) Publish(eventName string, eventData any) error {
	event, err := createEvent(eventName, eventData)
	if err != nil {
		return err
	}

	return eb.externalBus.Publish(event)
}

func (es EventSubscription) handle(busChannel <-chan Event, ctx context.Context) {
	for {
		select {
		case event := <-busChannel:
			go es.Handler.Handle(event.Data)
		case <-ctx.Done():
			return
		}
	}
}

func createEvent(eventName string, eventData any) (*Event, error) {
	dataEvent, err := serialize(eventData)
	if err != nil {
		return nil, err
	}

	return &Event{
		Type: eventName,
		Data: dataEvent,
	}, nil
}

func serialize(data any) ([]byte, error) {
	return json.Marshal(data)
}
