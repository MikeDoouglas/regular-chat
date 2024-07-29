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
	Room          *Room
	NameGenerator *NameGenerator
	Logger        *zap.SugaredLogger
}

func NewHandler(room *Room, nameGenerator *NameGenerator, logger *zap.SugaredLogger) *Handler {
	return &Handler{room, nameGenerator, logger}
}

func (h *Handler) HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Logger.Error("failed to upgrade HTTP to websocket", zap.Error(err))
		return
	}
	defer conn.Close()

	userId := uuid.New().String()
	name := h.NameGenerator.Generate()
	user := &User{Id: userId, Name: name, Conn: conn}

	h.Room.AddUser(user)
	defer func() {
		h.Room.RemoveUser(user)
		h.Logger.Infow("Disconnected", "user", user)
	}()

	if err := h.sendUserInfo(user); err != nil {
		h.Logger.Errorw("failed to send user info", zap.Error(err))
		h.sendError(err, user)
		return
	}

	h.Logger.Infow("New connection", "user_name", name, "user_id", userId)

	if err := h.listenNewMessages(conn); err != nil {
		h.Logger.Errorw("error while listening new messages", zap.Error(err))
		h.sendError(err, user)
	}
}

func (h *Handler) sendError(err error, user *User) {
	switch {
	case errors.Is(err, ErrDataSerialization):
		h.Room.NotifyError(user, &MessageJson{Type: MessageTypeError})
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

		h.Room.AddMessage(message)
		h.Room.NotifyUsers(message)

		h.Logger.Debugw("Received message", "message", message)
	}
}
