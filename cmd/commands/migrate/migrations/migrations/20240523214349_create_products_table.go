package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateProductsTable, downCreateProductsTable)
}

func upCreateProductsTable(ctx context.Context, tx *sql.Tx) error {
    var queries = `CREATE TABLE IF NOT EXISTS products (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        size VARCHAR(50) NOT NULL,
        sku VARCHAR(255) UNIQUE NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	_, err := tx.Exec(queries)
	if err != nil {
		return fmt.Errorf("error exec query. %w", err)
	}

	return nil
}

func downCreateProductsTable(ctx context.Context, tx *sql.Tx) error {
	var query = `DROP TABLE IF EXISTS products`
	_, err := tx.Exec(query)
	if err != nil {
		return fmt.Errorf("error drop products table. %w", err)
	}

	return nil
}
