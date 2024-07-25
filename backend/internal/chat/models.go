package chat

type MessageType string

const (
	MessageTypeError    MessageType = "ERROR"
	MessageTypeText     MessageType = "TEXT"
	MessageTypeUserInfo MessageType = "USER_INFO"
)

type MessageJson struct {
	Text     string      `json:"text"`
	UserId   string      `json:"user_id"`
	UserName string      `json:"user_name"`
	Type     MessageType `json:"type"`
}

type Message struct {
	Text   string
	UserId string
}
