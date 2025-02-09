// Package postgres implements postgres connection.
package postgres

import (
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Postgres -.
type Postgres struct {
	Builder squirrel.StatementBuilderType
	DB      *sqlx.DB
}

// New -.
func New(url string) (*Postgres, error) {
	db, openErr := sqlx.Open("pgx", url)
	if openErr != nil {
		log.Fatalf("sqlx.Open(): %v", openErr)
	}

	err := db.Ping()

	if err != nil {
		log.Fatalf("db.Ping(): %v", err)
	}

	pg := &Postgres{
		DB:      db,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	return pg, nil
}

// Close -.
func (p *Postgres) Close() {
	if p.DB != nil {
		err := p.DB.Close()
		if err != nil {
			log.Fatalf("db.Close(): %v", err)
		}
	}
}
