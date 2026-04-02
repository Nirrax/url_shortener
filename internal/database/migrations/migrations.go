package migrations

import (
	"database/sql"
	"fmt"

	"github.com/nirrax/url_shortener/internal/config"
	"github.com/pressly/goose/v3"
)

type GooseMigration struct {
	db              *sql.DB
	runUpMigrations bool
}

func NewGooseMigration(db *sql.DB, conf config.Config) *GooseMigration {
	return &GooseMigration{
		db:              db,
		runUpMigrations: conf.RunUpMigration,
	}
}

func (m *GooseMigration) Up() error {
	if !m.runUpMigrations {
		return nil
	}

	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	err = goose.Up(m.db, "internal/database/migrations/migrations/")
	if err != nil {
		return fmt.Errorf("goose up migration failed: %w", err)
	}

	return nil
}
