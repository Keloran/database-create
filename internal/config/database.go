package config

import (
	bugLog "github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v6"
)

type DatabaseDetails struct {
	Host     string
	Port     int
	UserName string
	Password string
}

type Database struct {
	Postgres DatabaseDetails `env:"DATABASE_POSTGRES"`
	Mongo    DatabaseDetails `env:"DATABASE_MONGO"`
	MySQL    DatabaseDetails `env:"DATABASE_MYSQL"`
}

func BuildDatabase(cfg *Config) error {
	databases := &Database{}

	// Mongo
	mongo, err := getMongoDetails()
	if err != nil {
		return bugLog.Error(err)
	}
	databases.Mongo = mongo

	// Postgres
	postgres, err := getPostgresDetails()
	if err != nil {
		return bugLog.Error(err)
	}
	databases.Postgres = postgres

	// Mysql
	mysql, err := getMySQLDetails()
	if err != nil {
		return bugLog.Error(err)

	}
	databases.MySQL = mysql

	cfg.Database = *databases
	return nil
}

func getMySQLDetails() (DatabaseDetails, error) {
	type MySQLDetails struct {
		Host     string `env:"DATABASE_MYSQL_HOST" envDefault:"cob.cobden.net"`
		Port     int    `env:"DATABASE_MYSQL_PORT" envDefault:"3306"`
		UserName string `env:"DATABASE_MYSQL_USERNAME"`
		Password string `env:"DATABASE_MYSQL_PASSWORD"`
	}

	mysql := &MySQLDetails{}
	if err := env.Parse(mysql); err != nil {
		return DatabaseDetails{}, bugLog.Error(err)
	}

	return DatabaseDetails{
		Host:     mysql.Host,
		Port:     mysql.Port,
		UserName: mysql.UserName,
		Password: mysql.Password,
	}, nil
}

func getPostgresDetails() (DatabaseDetails, error) {
	type PostgresDetails struct {
		Host     string `env:"DATABASE_POSTGRES_HOST" envDefault:"cob.cobden.net"`
		Port     int    `env:"DATABASE_POSTGRES_PORT" envDefault:"5432"`
		UserName string `env:"DATABASE_POSTGRES_USERNAME"`
		Password string `env:"DATABASE_POSTGRES_PASSWORD"`
	}

	postgres := &PostgresDetails{}
	if err := env.Parse(postgres); err != nil {
		return DatabaseDetails{}, bugLog.Error(err)
	}

	return DatabaseDetails{
		Host:     postgres.Host,
		Port:     postgres.Port,
		UserName: postgres.UserName,
		Password: postgres.Password,
	}, nil
}

func getMongoDetails() (DatabaseDetails, error) {
	type MongoDetails struct {
		Host     string `env:"DATABASE_MONGO_HOST" envDefault:"cob.cobden.net"`
		Port     int    `env:"DATABASE_MONGO_PORT" envDefault:"27017"`
		UserName string `env:"DATABASE_MONGO_USERNAME"`
		Password string `env:"DATABASE_MONGO_PASSWORD"`
	}

	mongo := &MongoDetails{}
	if err := env.Parse(mongo); err != nil {
		return DatabaseDetails{}, bugLog.Error(err)
	}

	return DatabaseDetails{
		Host:     mongo.Host,
		Port:     mongo.Port,
		UserName: mongo.UserName,
		Password: mongo.Password,
	}, nil
}
