package dbrepo

import (
	"database/sql"
	"myapp2/pkg/config"
	"myapp2/pkg/repository"
)

type postgresDb struct {
	App *config.AppConfig

	DB *sql.DB
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(con *sql.DB, a *config.AppConfig) repository.DatabaseRepo {

	return &postgresDb{
		App: a,
		DB:  con,
	}

}

func NewTestDBRepo(a *config.AppConfig) repository.DatabaseRepo {

	return &testDBRepo{
		App: a,
	}

}
