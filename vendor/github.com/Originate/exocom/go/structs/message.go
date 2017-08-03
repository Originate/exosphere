package structs

// Message defines the structure of websocket packets representing messages
type Message struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Payload    MessagePayload `json:"payload"`
	ResponseTo string         `json:"responseTo"`
	Sender     string         `json:"sender"`
	SessionID  string         `json:"sessionId"`
}

// MessagePayload defines the structure of Message.Payload
type MessagePayload interface{}
