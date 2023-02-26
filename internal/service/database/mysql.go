package database

import (
  "database/sql"
  "fmt"
  bugLog "github.com/bugfixes/go-bugfixes/logs"
  _ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
  Details
}

func NewMySQL(cfg Details) *MySQL {
  return &MySQL{
    Details: cfg,
  }
}

func (m *MySQL) DatabaseAlreadyExists(projectName string) (bool, error) {
  db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/", m.UserName, m.Password, m.Host, m.Port))
  if err != nil {
    return false, err
  }
  defer func() {
    if err := db.Close(); err != nil {
      bugLog.Local().Fatal(err)
    }
  }()

  rows, err := db.Query(fmt.Sprintf("SHOW DATABASES LIKE '%s'", projectName))
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

func (m *MySQL) Create(projectName string) error {
  db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/", m.UserName, m.Password, m.Host, m.Port))
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
