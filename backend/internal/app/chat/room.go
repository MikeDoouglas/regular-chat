package chat

import (
	"sync"

	"github.com/mikedoouglas/chat/internal/pkg/domain/model"
	"github.com/mikedoouglas/chat/internal/pkg/domain/user"
)

type Room struct {
	Users    map[string]*user.User
	Messages []*model.Message
	mu       sync.Mutex
}

func NewRoom() *Room {
	return &Room{
		Users: make(map[string]*user.User),
	}
}

func (r *Room) AddUser(user *user.User) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Users[user.Id] = user
}

func (r *Room) RemoveUser(user *user.User) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.Users, user.Id)
}

func (r *Room) AddMessage(message *model.MessageJSON) {
	r.Messages = append(r.Messages, &model.Message{
		UserId: message.UserId,
		Text:   message.Text,
	})
}

func (r *Room) NotifyUsers(message *model.MessageJSON) {
	r.mu.Lock()
	usersCopy := make([]*user.User, 0, len(r.Users))
	for _, u := range r.Users {
		usersCopy = append(usersCopy, u)
	}
	r.mu.Unlock()

	for _, u := range usersCopy {
		if message.UserId != u.Id {
			go u.SendMessage(message)
		}
	}
}

func (r *Room) NotifyError(user *user.User, message *model.MessageJSON) {
	user.SendMessage(message)
}
