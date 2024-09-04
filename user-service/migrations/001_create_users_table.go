package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(UP_001, Down_001)
}

func UP_001(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    phone VARCHAR(100),
    city VARCHAR(100),
    password VARCHAR(255),
    role VARCHAR(100),
    avatar_url VARCHAR(255),
    created_at bigint,
    deleted_At bigint
);`)

	if err != nil {
		return err
	}
	return nil
}

func Down_001(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE users;`)
	if err != nil {
		return err
	}
	return nil
}
