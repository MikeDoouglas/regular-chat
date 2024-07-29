package chat

import "sync"

type Room struct {
	Users    []*User
	Messages []*Message
	mutex    *sync.Mutex
}

func NewRoom(mutex *sync.Mutex) *Room {
	return &Room{mutex: mutex}
}

func (r *Room) AddUser(user *User) {
	r.Users = append(r.Users, user)
}

func (r *Room) RemoveUser(user *User) {
	r.mutex.Lock()
	index := r.findUserIndex(user.Id)
	r.Users = append(r.Users[:index], r.Users[index+1:]...)
	r.mutex.Unlock()
}

func (r *Room) AddMessage(message *MessageJson) {
	r.Messages = append(r.Messages, &Message{
		UserId: message.UserId,
		Text:   message.Text,
	})
}

func (r *Room) NotifyUsers(message *MessageJson) {
	for _, u := range r.Users {
		if message.UserId != u.Id {
			u.SendMessage(message, r.mutex)
		}
	}
}

func (r *Room) NotifyError(user *User, message *MessageJson) {
	user.SendMessage(message, r.mutex)
}

func (r *Room) findUserIndex(userId string) int {
	for index, u := range r.Users {
		if u.Id == userId {
			return index
		}
	}

	return -1
}
