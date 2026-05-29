package event_bus

type Event struct {
	Name string
	Data any
}

const (
	EventClickCategory         = "click_category"
	EventClickArticle          = "click_article"
	EventCreateUserArticle     = "create_article"
	EventUpdateUserArticle     = "update_article"
	EventSoftDeleteUserArticle = "soft_delete_article"
	EventHardDeleteUserArticle = "hard_delete_article"
	EventRecoveryUserArticle   = "recovery_article"
)

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
