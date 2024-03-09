package user

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

func (r *Repository) CreateUser(credentials entities.AuthCredentials) error {
	return r.db.Insert(fmt.Sprintf("INSERT INTO users(username, password) VALUES('%s', '%s')", credentials.Login, credentials.Password), struct{}{})
}

func (r *Repository) GetUser(username string) (entities.AuthCredentials, error) {
	query := fmt.Sprintf("SELECT password FROM users WHERE username = '%s'", username)

	data, err := r.db.Get(query)
	if err != nil {
		return entities.AuthCredentials{}, err
	}

	rows, ok := data.(*sqlx.Rows)
	if !ok {
		return entities.AuthCredentials{}, ErrIncorrectTypeRows
	}

	rows.Next()

	var user models.User

	err = rows.StructScan(&user)
	if err != nil {
		return entities.AuthCredentials{}, err
	}

	return mapper.MakeEntityAuthCredentials(username, user.Password), err
}
