package dbrepository

import (
	"database/sql"

	"github.com/bharatsabne/bookings/Internal/config"
	"github.com/bharatsabne/bookings/Internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.Databaserepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
