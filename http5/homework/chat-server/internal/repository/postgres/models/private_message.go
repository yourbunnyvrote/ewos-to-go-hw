package models

type PrivateMessage struct {
	Receiver string `db:"receiver"`
	Sender   string `db:"sender"`
	Message  string `db:"message"`
}
