package repository

import chatutil "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server"

type Authorization interface {
	CreateUser(user chatutil.User) (string, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *chatutil.ChatDB) *Repository {
	return &Repository{Authorization: NewAuthDB(db)}
}
