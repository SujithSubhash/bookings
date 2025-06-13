package dbrepo

import (
	"database/sql"

	"github.com/SujithSubhash/bookings/internal/config"
	"github.com/SujithSubhash/bookings/repository"
)

type PostgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &PostgresDBRepo{
		App: a,
		DB:  conn,
	}
}
