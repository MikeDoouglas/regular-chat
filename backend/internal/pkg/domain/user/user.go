package user

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/mikedoouglas/chat/internal/pkg/domain/model"
	"go.uber.org/zap"
)

type User struct {
	Id     string
	Name   string
	conn   *websocket.Conn
	logger *zap.SugaredLogger
}

func New(id, name string, conn *websocket.Conn, logger *zap.SugaredLogger) *User {
	return &User{
		Id:     id,
		Name:   name,
		conn:   conn,
		logger: logger,
	}
}

func (u *User) SendMessage(m *model.MessageJSON) {
	message, err := json.Marshal(m)
	if err != nil {
		u.logger.Errorw("failed to marshal message when sending to user",
			"message", message, zap.Error(err))
	}

	if err = u.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		u.logger.Errorw("failed to send message to user",
			"message", message, zap.Error(err))
	}
}

func (u *User) WriteMessage(messageType int, data []byte) error {
	return u.conn.WriteMessage(messageType, data)
}
