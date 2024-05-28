package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateUsersTable, downCreateUsersTable)
}

func upCreateUsersTable(ctx context.Context, tx *sql.Tx) error {
	var queries = `CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        email VARCHAR(100) NOT NULL,
        phone VARCHAR(20),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`

	_, err := tx.Exec(queries)
	if err != nil {
		return fmt.Errorf("error exec query. %w", err)
	}

	return nil
}

func downCreateUsersTable(ctx context.Context, tx *sql.Tx) error {
	var query = `DROP TABLE IF EXISTS users`
	_, err := tx.Exec(query)
	if err != nil {
		return fmt.Errorf("error drop users table. %w", err)
	}

	return nil
}
