package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("postgres", "postgresql://sample_ojdt_user:BNfKxdWCRRUtG12abmk3ILyk74a26dyY@dpg-cq1cv608fa8c739osj90-a.singapore-postgres.render.com/sample_ojdt")
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
