package config

import (
  bugLog "github.com/bugfixes/go-bugfixes/logs"
  "github.com/caarlos0/env/v6"
)

type Config struct {
  Local
  Database
  Vault
}

func Build() (*Config, error) {
  cfg := &Config{}

  if err := env.Parse(cfg); err != nil {
    return nil, bugLog.Error(err)
  }

  if err := BuildLocal(cfg); err != nil {
    return nil, bugLog.Error(err)
  }

  if err := BuildDatabase(cfg); err != nil {
    return nil, bugLog.Error(err)
  }

  if err := BuildVault(cfg); err != nil {
    return nil, bugLog.Error(err)
  }

  return cfg, nil
}
