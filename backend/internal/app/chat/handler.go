package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mikedoouglas/chat/internal/pkg/domain/model"
	"github.com/mikedoouglas/chat/internal/pkg/domain/user"
	namegenerator "github.com/mikedoouglas/chat/internal/pkg/generator"
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
	nameGenerator *namegenerator.NameGenerator
	logger        *zap.SugaredLogger
}

func NewHandler(room *Room, nameGenerator *namegenerator.NameGenerator, logger *zap.SugaredLogger) *Handler {
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
	user := user.New(userId, name, conn, h.logger)

	h.room.AddUser(user)
	defer func() {
		h.room.RemoveUser(user)
		h.logger.Infow("Disconnected", "user_id", user.Id, "user_name", user.Name)
	}()

	if err := h.sendUserInfo(user); err != nil {
		h.logger.Errorw("failed to send user info", zap.Error(err))
		h.sendError(err, user)
		return
	}

	h.logger.Infow("New connection", "user_name", name, "user_id", userId)

	if err := h.listenNewMessages(conn); err != nil {
		var closeErr *websocket.CloseError
		if errors.As(err, &closeErr) &&
			(closeErr.Code == websocket.CloseGoingAway || closeErr.Code == websocket.CloseNormalClosure) {
			h.logger.Infow("client disconnected", zap.Error(err))
		} else {
			h.logger.Errorw("unexpected error while listening new messages", zap.Error(err))
			h.sendError(err, user)
		}
	}
}

func (h *Handler) sendError(err error, user *user.User) {
	switch {
	case errors.Is(err, ErrDataSerialization):
		h.room.NotifyError(user, &model.MessageJSON{Type: model.MessageTypeError})
	default:
		msg := websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "internal server error")
		user.WriteMessage(websocket.CloseMessage, msg)
	}

}

func (h *Handler) sendUserInfo(user *user.User) error {
	m := &model.MessageJSON{UserId: user.Id, UserName: user.Name, Type: model.MessageTypeUserInfo}
	message, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal user info: %w", ErrDataSerialization)
	}

	return user.WriteMessage(websocket.TextMessage, message)
}

func (h *Handler) listenNewMessages(conn *websocket.Conn) error {
	for {
		_, m, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		var message *model.MessageJSON
		if err := json.Unmarshal(m, &message); err != nil {
			return fmt.Errorf("failed to unmarshal received message: %w", ErrDataSerialization)
		}

		h.room.AddMessage(message)
		h.room.NotifyUsers(message)

		h.logger.Debugw("Received message", "message", message)
	}
}
