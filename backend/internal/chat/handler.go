package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var ErrDataSerialization = errors.New("failed to serialize data")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	room          *Room
	nameGenerator *NameGenerator
	logger        *zap.SugaredLogger
}

func NewHandler(room *Room, nameGenerator *NameGenerator, logger *zap.SugaredLogger) *Handler {
	return &Handler{room, nameGenerator, logger}
}

func (h *Handler) HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("failed to upgrade HTTP to websocket", zap.Error(err))
		return
	}
	defer conn.Close()

	userId := uuid.New().String()
	name := h.nameGenerator.Generate()
	user := &User{Id: userId, Name: name, Conn: conn}

	h.room.AddUser(user)
	defer func() {
		h.room.RemoveUser(user)
		h.logger.Infow("Disconnected", "user", user)
	}()

	if err := h.sendUserInfo(user); err != nil {
		h.logger.Errorw("failed to send user info", zap.Error(err))
		h.sendError(err, user)
		return
	}

	h.logger.Infow("New connection", "user_name", name, "user_id", userId)

	if err := h.listenNewMessages(conn); err != nil {
		h.logger.Errorw("error while listening new messages", zap.Error(err))
		h.sendError(err, user)
	}
}

func (h *Handler) sendError(err error, user *User) {
	switch {
	case errors.Is(err, ErrDataSerialization):
		h.room.NotifyError(user, &MessageJson{Type: MessageTypeError})
	default:
		user.Conn.WriteMessage(websocket.CloseInternalServerErr, nil)
	}

}

func (h *Handler) sendUserInfo(user *User) error {
	m := &MessageJson{UserId: user.Id, UserName: user.Name, Type: MessageTypeUserInfo}
	message, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal user info: %w", ErrDataSerialization)
	}

	return user.Conn.WriteMessage(websocket.TextMessage, message)
}

func (h *Handler) listenNewMessages(conn *websocket.Conn) error {
	for {
		_, m, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		var message *MessageJson
		if err := json.Unmarshal(m, &message); err != nil {
			return fmt.Errorf("failed to unmarshal received message: %w", ErrDataSerialization)
		}

		h.room.AddMessage(message)
		h.room.NotifyUsers(message)

		h.logger.Debugw("Received message", "message", message)
	}
}
