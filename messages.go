package jarvis

// Message codes
const (
	ActionRegister = "register"
)

// ClientMessage is the format that server expects from clients
type ClientMessage struct {
	ID     string `json:"id"`
	Action string `json:"action"`
}

// ServerMessage is the format that clients expect from the server
type ServerMessage struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}
