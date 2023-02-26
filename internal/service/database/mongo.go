package database

import (
	"context"
	"fmt"
	bugLog "github.com/bugfixes/go-bugfixes/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Mongo struct {
	Details Details
}

func NewMongo(details Details) *Mongo {
	return &Mongo{
		Details: details,
	}
}

func (m *Mongo) DatabaseAlreadyExists(projectName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", m.Details.Host, m.Details.Port)),
		options.Client().SetAuth(options.Credential{
			Username: m.Details.UserName,
			Password: m.Details.Password,
		}))
	if err != nil {
		return false, err
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			bugLog.Local().Fatal(err)
		}
	}()
	dbs, err := client.ListDatabaseNames(ctx, bson.M{"name": projectName})
	if err != nil {
		return false, err
	}

	if len(dbs) > 0 {
		return true, nil
	}
	return false, nil
}

func (m *Mongo) Create(projectName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", m.Details.Host, m.Details.Port)),
		options.Client().SetAuth(options.Credential{
			Username: m.Details.UserName,
			Password: m.Details.Password,
		}))
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			bugLog.Local().Fatal(err)
		}
	}()
	db := client.Database(projectName)
	if err := db.CreateCollection(ctx, projectName); err != nil {
		return err
	}

	return nil
}
