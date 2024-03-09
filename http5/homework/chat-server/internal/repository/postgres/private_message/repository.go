package private_message

import (
	"errors"
	"fmt"

	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	"github.com/ew0s/ewos-to-go-hw/internal/repository/postgres/models"
	"github.com/jmoiron/sqlx"
)

var ErrIncorrectTypeRows = errors.New("rows are not of type *sqlx.Rows")

type Postgres interface {
	Insert(query string, _ interface{}) error
	Get(query string) (interface{}, error)
}

type Repository struct {
	db Postgres
}

func NewRepository(db Postgres) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) SendPrivateMessage(receiver string, msg entities.Message) error {
	return r.db.Insert(fmt.Sprintf("INSERT INTO private_messages(receiver, sender, message) VALUES('%s', '%s', '%s')", receiver, msg.Username, msg.Content), struct{}{})
}

func (r *Repository) GetPrivateChat(chat entities.ChatMetadata) ([]entities.Message, error) {
	data, err := r.db.Get(fmt.Sprintf("SELECT sender, message FROM private_messages WHERE receiver = '%s' AND sender = '%s' UNION SELECT sender, message FROM private_messages WHERE receiver = '%s' AND sender = '%s' ORDER BY id", chat.Username1, chat.Username2, chat.Username2, chat.Username1))
	if err != nil {
		return nil, err
	}

	var (
		msg       models.PrivateMessage
		resultMsg entities.Message
	)

	result := make([]entities.Message, 0)

	rows, ok := data.(*sqlx.Rows)
	if !ok {
		return nil, ErrIncorrectTypeRows
	}

	for rows.Next() {
		err = rows.StructScan(&msg)
		if err != nil {
			return nil, err
		}

		resultMsg.Username = msg.Sender
		resultMsg.Content = msg.Message
		result = append(result, resultMsg)
	}

	return result, nil
}

func (r *Repository) GetUserList(username string) ([]string, error) {
	subQuery1 := fmt.Sprintf("SELECT receiver as username FROM private_messages WHERE sender = '%s'", username)
	subQuery2 := fmt.Sprintf("SELECT sender as username FROM private_messages WHERE receiver = '%s'", username)

	query := fmt.Sprintf("SELECT DISTINCT username FROM ((%s) UNION (%s)) as all_values", subQuery1, subQuery2)

	data, err := r.db.Get(query)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)

	var usernameLists string

	rows, ok := data.(*sqlx.Rows)
	if !ok {
		return nil, ErrIncorrectTypeRows
	}

	for rows.Next() {
		err = rows.Scan(&usernameLists)
		if err != nil {
			return nil, err
		}

		result = append(result, usernameLists)
	}

	return result, nil
}
