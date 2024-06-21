package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type databaseI interface {
	CreateUrl(Url) error
	DeleteUrl(int) error
	GetUrlByID(int) (*Url, error)
	GetUrlByShortUrl(string) (*Url, error)
}

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB() (*PostgresDB, error) {
	connStr := "user=postgres dbname=postgresdb password=postgrespass sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresDB{
		db: db,
	}, nil
}

func (db *PostgresDB) Init() error {
	return db.CreateTables()
}

func (db *PostgresDB) CreateTables() error {
	// 2083 - value from stackOverflow (practical limit of the http protocole)
	// 8 - length of the shortened url
	query := `CREATE TABLE IF NOT EXISTS url (
		id INT PRIMARY KEY,
		long_url VARCHAR(2083), 
		short_url VARCHAR(8),
		created_at TIMESTAMP
	) `

	_, err := db.db.Exec(query)
	return err
}

func (db *PostgresDB) CreateUrl(url Url) error {
	query := fmt.Sprintf(
		`INSERT INT url
		(long_url, short_url, created_at) 
		VALUES (%v , %v, %v)`, url.LongUrl, url.ShortUrl, url.CreatedAt,
	)

	_, err := db.db.Exec(query)
	return err
}

func (db *PostgresDB) DeleteUrl(id int) error {
	query := fmt.Sprintf("DELETE FROM url WHERE id=%v", id)

	_, err := db.db.Query(query)
	return err
}

func (db *PostgresDB) GetUrlByID(id int) (*Url, error) {
	query := fmt.Sprintf(`SELECT * FROM url WHERE id=%v`, id)

	row := db.db.QueryRow(query)

	object, err := scanIntoUrl(row)

	if err != nil {
		return nil, err
	}

	return object, nil
}

func (db *PostgresDB) GetUrlByShortUrl(url string) (*Url, error) {
	query := fmt.Sprintf(`SELECT * FROM url WHERE short_url=%v`, url)

	row := db.db.QueryRow(query)

	object, err := scanIntoUrl(row)

	if err != nil {
		return nil, err
	}

	return object, nil
}

func scanIntoUrl(rows *sql.Row) (*Url, error) {
	url := new(Url)
	err := rows.Scan(
		&url.ID,
		&url.LongUrl,
		&url.ShortUrl,
		&url.CreatedAt,
	)

	return url, err
}
