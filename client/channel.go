package client

// Handler defines handling of channel events
type Handler interface {
	HandleEvent(Event)
}

// HandlerFunc defines functions that handle events
type HandlerFunc func(Event)

// HandleEvent handles an event with a HandlerFunc
func (f HandlerFunc) HandleEvent(e Event) {
	f(e)
}

// Channel defines a connection to a websocket
type Channel interface {
	Handler

	UpdateClientState(connected bool)

	SetActive(active bool)

	Subscribe()
	Unsubscribe()

	Bind(event string, h Handler)
	Unbind(event string, h Handler)

	BindAll(h Handler)
	UnbindAll(h Handler)

	BindFunc(event string, h func(Event))
	UnbindFunc(event string, h func(Event))

	BindAllFunc(h func(Event))
	UnbindAllFunc(h func(Event))
}
