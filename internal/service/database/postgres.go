package database

import (
	"database/sql"
	"fmt"
	bugLog "github.com/bugfixes/go-bugfixes/logs"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Details
}

func NewPostgres(cfg Details) *Postgres {
	return &Postgres{
		Details: cfg,
	}
}

func (p *Postgres) DatabaseAlreadyExists(projectName string) (bool, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d?sslmode=disable", p.UserName, p.Password, p.Host, p.Port))
	if err != nil {
		return false, err
	}
	defer func() {
		if err := db.Close(); err != nil {
			bugLog.Local().Fatal(err)
		}
	}()
	rows, err := db.Query(fmt.Sprintf("SELECT datname FROM pg_database WHERE datname = '%s'", projectName))
	if err != nil {
		return false, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			bugLog.Local().Fatal(err)
		}
	}()
	if rows.Next() {
		return true, nil
	}

	return false, nil
}

func (p *Postgres) Create(projectName string) error {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d?sslmode=disable", p.UserName, p.Password, p.Host, p.Port))
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			bugLog.Local().Fatal(err)
		}
	}()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", projectName))
	if err != nil {
		return err
	}

	return nil
}
