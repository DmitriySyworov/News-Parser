package event_bus

type Event struct {
	Name string
	Data any
}

type EventBus struct {
	Bus chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{
		Bus: make(chan Event),
	}
}
func (e *EventBus) Publisher(event *Event) {
	e.Bus <- *event
}
func (e *EventBus) Subscriber() <-chan Event {
	return e.Bus
}
