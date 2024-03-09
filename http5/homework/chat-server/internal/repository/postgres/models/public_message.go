package models

type PublicMessage struct {
	Username string `db:"sender"`
	Message  string `db:"message"`
}
