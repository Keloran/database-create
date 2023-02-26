package config

import (
  bugLog "github.com/bugfixes/go-bugfixes/logs"
  "github.com/caarlos0/env/v6"
)

type Vault struct {
  Token   string `env:"VAULT_TOKEN"`
  Address string `env:"VAULT_ADDRESS"`
}

func BuildVault(cfg *Config) error {
  vault := &Vault{}
  if err := env.Parse(vault); err != nil {
    return bugLog.Error(err)
  }
  cfg.Vault = *vault
  return nil
}
