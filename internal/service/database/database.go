package database

import "github.com/keloran/database-create/internal/config"

type Database struct {
	Config *config.Config
}

type Details struct {
	Host     string
	Port     int
	UserName string
	Password string
}

func NewDatabase(cfg *config.Config) *Database {
	return &Database{
		Config: cfg,
	}
}

type System interface {
	DatabaseAlreadyExists(projectName string) (bool, error)
	Create(projectName string) error
}

func (d *Database) FetchSystem(system string) System {
	switch system {
	case "mysql":
		return NewMySQL(Details{
			Host:     d.Config.Database.MySQL.Host,
			Port:     d.Config.Database.MySQL.Port,
			UserName: d.Config.Database.MySQL.UserName,
			Password: d.Config.Database.MySQL.Password,
		})
	case "postgres":
		return NewPostgres(Details{
			Host:     d.Config.Database.Postgres.Host,
			Port:     d.Config.Database.Postgres.Port,
			UserName: d.Config.Database.Postgres.UserName,
			Password: d.Config.Database.Postgres.Password,
		})
	case "mongo":
		return NewMongo(Details{
			Host:     d.Config.Database.Mongo.Host,
			Port:     d.Config.Database.Mongo.Port,
			UserName: d.Config.Database.Mongo.UserName,
			Password: d.Config.Database.Mongo.Password,
		})
	}

	return nil
}
