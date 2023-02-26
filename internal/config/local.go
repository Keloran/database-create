package config

import (
  bugLog "github.com/bugfixes/go-bugfixes/logs"
  "github.com/caarlos0/env/v6"
)

type Local struct {
  KeepLocal   bool `env:"KEEP_LOCAL" envDefault:"false"`
  Development bool `env:"DEVELOPMENT" envDefault:"false"`
}

func BuildLocal(cfg *Config) error {
  local := &Local{}

  if err := env.Parse(local); err != nil {
    return bugLog.Error(err)
  }

  cfg.Local = *local
  return nil
}
