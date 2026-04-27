package event_bus

type Event struct {
	Name string
	Data any
}
type Bus struct {
	EventBus chan Event
}

func (b *Bus) Publisher(event Event) {
	b.EventBus <- event
}
func (b *Bus) Subscriber() chan Event {
	return b.EventBus
}
