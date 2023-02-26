package vault

import (
	"fmt"
	bugLog "github.com/bugfixes/go-bugfixes/logs"
	"github.com/hashicorp/vault/api"
	"github.com/keloran/database-create/internal/config"
)

type Vault struct {
	Token   string
	Address string
	Config  *config.Config
}

func NewVault(cfg *config.Config) *Vault {
	return &Vault{
		Token:   cfg.Vault.Token,
		Address: cfg.Vault.Address,
		Config:  cfg,
	}
}

func (v *Vault) DatabaseAlreadyExists(projectName string) (bool, error) {
	client, err := v.getClient()
	if err != nil {
		return false, bugLog.Error(err)
	}

	data, err := client.Logical().Read(fmt.Sprintf("database/creds/%s-database-role", projectName))
	if err != nil {
		if resp, ok := err.(*api.ResponseError); ok {
			if resp.StatusCode != 400 {
				return false, bugLog.Error(err)
			}
		}
	}

	if data != nil {
		return true, nil
	}

	return false, nil
}

func (v *Vault) getClient() (*api.Client, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, bugLog.Error(err)
	}

	client.SetToken(v.Token)
	if err := client.SetAddress(v.Address); err != nil {
		return nil, bugLog.Error(err)
	}

	return client, nil
}

func (v *Vault) CreateMySQL(projectName string) error {
	client, err := v.getClient()
	if err != nil {
		return bugLog.Error(err)
	}

	// Create connection
	_, err = client.Logical().Write(fmt.Sprintf("database/config/%s-mysql", projectName), map[string]interface{}{
		"plugin_name":   "mysql-database-plugin",
		"allowed_roles": fmt.Sprintf("%s-database-role", projectName),
		"connection_url": fmt.Sprintf(
			"{{username}}:{{password}}@tcp(%s:%d)/",
			v.Config.Database.MySQL.Host,
			v.Config.Database.MySQL.Port),
		"username": v.Config.Database.MySQL.UserName,
		"password": v.Config.Database.MySQL.Password,
	})
	if err != nil {
		return bugLog.Error(err)
	}

	// Create Role
	_, err = client.Logical().Write(fmt.Sprintf("database/roles/%s-database-role", projectName), map[string]interface{}{
		"db_name":             fmt.Sprintf("%s-mysql", projectName),
		"creation_statements": fmt.Sprintf("CREATE USER '{{name}}'@'%%' IDENTIFIED BY '{{password}}';GRANT ALL PRIVILEGES ON %s.* TO '{{name}}'@'%%';", projectName),
		"default_ttl":         "1h",
		"max_ttl":             "24h",
	})
	if err != nil {
		return bugLog.Error(err)
	}

	return nil
}

func (v *Vault) CreatePostgres(projectName string) error {
	client, err := v.getClient()
	if err != nil {
		return bugLog.Error(err)
	}

	// Create connection
	_, err = client.Logical().Write(fmt.Sprintf("database/config/%s-postgres", projectName), map[string]interface{}{
		"plugin_name":   "postgresql-database-plugin",
		"allowed_roles": fmt.Sprintf("%s-database-role", projectName),
		"connection_url": fmt.Sprintf(
			"postgres://{{username}}:{{password}}@%s:%d/%s?sslmode=disable",
			v.Config.Database.Postgres.Host,
			v.Config.Database.Postgres.Port,
			projectName),
		"username": v.Config.Database.Postgres.UserName,
		"password": v.Config.Database.Postgres.Password,
	})
	if err != nil {
		return bugLog.Error(err)
	}

	// Create Role
	_, err = client.Logical().Write(fmt.Sprintf("database/roles/%s-database-role", projectName), map[string]interface{}{
		"db_name":             fmt.Sprintf("%s-postgres", projectName),
		"creation_statements": fmt.Sprintf("CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; GRANT ALL PRIVILEGES ON DATABASE %s TO \"{{name}}\";", projectName),
		"default_ttl":         "1h",
		"max_ttl":             "24h",
	})
	if err != nil {
		return bugLog.Error(err)
	}

	return nil
}

func (v *Vault) CreateMongo(projectName string) error {
	client, err := v.getClient()
	if err != nil {
		return bugLog.Error(err)
	}

	// Create connection
	_, err = client.Logical().Write(fmt.Sprintf("database/config/%s-mongo", projectName), map[string]interface{}{
		"plugin_name":   "mongodb-database-plugin",
		"allowed_roles": fmt.Sprintf("%s-database-role", projectName),
		"connection_url": fmt.Sprintf(
			"mongodb://{{username}}:{{password}}@%s:%d/admin?tls=false",
			v.Config.Database.Mongo.Host,
			v.Config.Database.Mongo.Port),
		"username": v.Config.Database.Mongo.UserName,
		"password": v.Config.Database.Mongo.Password,
	})
	if err != nil {
		return bugLog.Error(err)
	}

	// Create Role
	_, err = client.Logical().Write(fmt.Sprintf("database/roles/%s-database-role", projectName), map[string]interface{}{
		"db_name":             fmt.Sprintf("%s-mongo", projectName),
		"creation_statements": fmt.Sprintf(`{"db": "admin", "roles": [{"role":"readWrite"}, {"role": "readWrite", "db": "%s"]}`, projectName),
		"default_ttl":         "1h",
		"max_ttl":             "24h",
	})
	if err != nil {
		return bugLog.Error(err)
	}

	return nil
}
