package chat

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type User struct {
	Id     string
	Name   string
	Conn   *websocket.Conn
	logger *zap.SugaredLogger
}

func (u *User) SendMessage(m *MessageJson, mu *sync.Mutex) {
	message, err := json.Marshal(m)
	if err != nil {
		u.logger.Errorw("failed to marshal message when sending to user",
			"message", message, zap.Error(err))
	}

	mu.Lock()
	if err = u.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
		u.logger.Errorw("failed to send message to user",
			"message", message, zap.Error(err))
	}
	mu.Unlock()
}
