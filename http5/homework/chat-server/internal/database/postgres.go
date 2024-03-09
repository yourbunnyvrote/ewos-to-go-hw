package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PostgresDB struct {
	db *sqlx.DB
}

func NewPostgresDB() *PostgresDB {
	db, err := connectDB()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return &PostgresDB{
		db: db,
	}
}

func connectDB() (*sqlx.DB, error) {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "postgres"
		dbname   = "postgres"
	)

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (p *PostgresDB) Insert(query string, _ interface{}) error {
	_, err := p.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) Get(query string) (interface{}, error) {
	rows, err := p.db.Queryx(query)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
