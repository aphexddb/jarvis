package client

// Event defines handlers for events that occur
type Event interface {
	GetEvent() string
	SetEvent(string)
	GetChannel() string
	SetChannel(string)

	GetData() interface{}
	SetData(interface{})
	GetDataString() string
	SetDataString(string)
}
