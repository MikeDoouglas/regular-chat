package chat

type Room struct {
	Users    []*User
	Messages []*Message
}

func (r *Room) AddUser(user *User) {
	r.Users = append(r.Users, user)
}

func (r *Room) RemoveUser(user *User) {
	index := r.findUserIndex(user.Id)
	r.Users = append(r.Users[:index], r.Users[index+1:]...)
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
			u.SendMessage(message)
		}
	}
}

func (r *Room) findUserIndex(userId string) int {
	for index, u := range r.Users {
		if u.Id == userId {
			return index
		}
	}

	return -1
}
