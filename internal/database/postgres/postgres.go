package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nirrax/url_shortener/internal/config"
)

type PostgresDB struct {
	*sqlx.DB
}

func NewPostgresDB(conf config.PostgresConfig) (*PostgresDB, error) {
	dsn := generateDSN(conf)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	return &PostgresDB{
		db,
	}, nil
}

func (p *PostgresDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return p.DB.Exec(query, args...)
}

func (p *PostgresDB) FetchOne(dest any, query string, args ...any) error {
	return p.DB.Get(dest, query, args...)
}

func (p *PostgresDB) FetchMany(dest any, query string, args ...any) error {
	return p.DB.Select(dest, query, args...)
}

func generateDSN(conf config.PostgresConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", conf.DBhost, conf.DBport, conf.DBuser, conf.DBpassword, conf.DBname)
}
