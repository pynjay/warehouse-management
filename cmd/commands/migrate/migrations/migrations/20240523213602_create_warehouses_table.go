package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateWarehousesTable, downCreateWarehousesTable)
}

func upCreateWarehousesTable(ctx context.Context, tx *sql.Tx) error {
	var queries = `CREATE TABLE IF NOT EXISTS warehouses (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        is_available BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`

	_, err := tx.Exec(queries)
	if err != nil {
		return fmt.Errorf("error exec query. %w", err)
	}

	return nil
}

func downCreateWarehousesTable(ctx context.Context, tx *sql.Tx) error {
	var query = `DROP TABLE IF EXISTS warehouses`
	_, err := tx.Exec(query)
	if err != nil {
		return fmt.Errorf("error drop warehouses table. %w", err)
	}

	return nil
}
