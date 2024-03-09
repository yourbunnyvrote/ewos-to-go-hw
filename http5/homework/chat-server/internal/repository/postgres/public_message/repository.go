package public_message

import (
	"errors"
	"fmt"

	"github.com/ew0s/ewos-to-go-hw/internal/api/mapper"
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

func (r *Repository) SendPublicMessage(msg entities.Message) error {
	query := fmt.Sprintf("INSERT INTO public_messages(sender, message) VALUES('%s', '%s')", msg.Username, msg.Content)
	return r.db.Insert(query, struct{}{})
}

func (r *Repository) GetPublicChat() ([]entities.Message, error) {
	data, err := r.db.Get("SELECT sender, message FROM public_messages ORDER BY id")
	if err != nil {
		return nil, err
	}

	result := make([]entities.Message, 0)

	var msg models.PublicMessage

	rows, ok := data.(*sqlx.Rows)
	if !ok {
		return nil, ErrIncorrectTypeRows
	}

	for rows.Next() {
		err = rows.StructScan(&msg)
		if err != nil {
			return nil, err
		}

		result = append(result, mapper.MakeEntityMessage(msg.Username, msg.Message))
	}

	return result, nil
}
