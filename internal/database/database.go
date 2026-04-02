package database

import "database/sql"

type DBclient interface {
	// Exec run sql query
	// FetchOne fetch single row into struct
	// FetchMany fetch multiple rows into slice of structs
	Exec(query string, args ...any) (sql.Result, error)
	FetchOne(dest any, query string, args ...any) error
	FetchMany(dest any, query string, args ...any) error
}
