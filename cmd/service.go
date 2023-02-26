package main

import (
  bugLog "github.com/bugfixes/go-bugfixes/logs"
  "github.com/joho/godotenv"
  "github.com/keloran/database-create/internal/config"
  "github.com/keloran/database-create/internal/service/database"
  "github.com/keloran/database-create/internal/service/vault"
  "github.com/spf13/cobra"
  "os"
  "strings"
)

var databaseTypes = []string{
  "mysql",
  "postgres",
  "mongo",
}

func main() {
  var (
    project      string
    databaseType string
  )

  if err := godotenv.Load(); err != nil {
    bugLog.Local().Info(err)
    os.Exit(1)
  }

  cfg, err := config.Build()
  if err != nil {
    bugLog.Local().Info(err)
    os.Exit(1)
  }

  if !cfg.Development {
    var cobraCommand = &cobra.Command{
      Use:   "database",
      Short: "Create the database and its vault account e.g database <project> <database type>",
      Args:  cobra.RangeArgs(1, 2),
      Run: func(cmd *cobra.Command, args []string) {
        project = args[0]

        if len(args) == 2 {
          if databaseValid(args[1]) {
            databaseType = args[1]
          } else {
            bugLog.Local().Info("Invalid database type\n Valid types are: " + strings.Join(databaseTypes, ",") + "\n")
            os.Exit(1)
          }
        }
      },
    }

    if err := cobraCommand.Execute(); err != nil {
      bugLog.Local().Info(err)
      os.Exit(1)
    }
  } else {
    project = os.Getenv("TEST_PROJECT")
    databaseType = os.Getenv("TEST_ENGINE")
  }

  if project == "" {
    bugLog.Local().Info("No project name given")
    os.Exit(1)
  }

  if err := createDatabase(cfg, databaseType, project); err != nil {
    bugLog.Local().Info(err)
    os.Exit(1)
  }

  if err := createVault(cfg, databaseType, project); err != nil {
    bugLog.Local().Info(err)
    os.Exit(1)
  }
}

func databaseValid(database string) bool {
  for _, databaseType := range databaseTypes {
    if database == databaseType {
      return true
    }
  }

  return false
}

func createVault(cfg *config.Config, databaseType, projectName string) error {
  v := vault.NewVault(cfg)

  switch databaseType {
  case "mysql":
    if err := v.CreateMySQL(projectName); err != nil {
      return bugLog.Error(err)
    }

  case "postgres":
    if err := v.CreatePostgres(projectName); err != nil {
      return bugLog.Error(err)
    }

  case "mongo":
    if err := v.CreateMongo(projectName); err != nil {
      return bugLog.Error(err)
    }
  }
  return nil
}

func createDatabase(cfg *config.Config, databaseType, projectName string) error {
  d := database.NewDatabase(cfg)
  s := d.FetchSystem(databaseType)
  exists, err := s.DatabaseAlreadyExists(projectName)
  if err != nil {
    return bugLog.Error(err)
  }

  if !exists {
    return s.Create(projectName)
  }

  return nil
}
