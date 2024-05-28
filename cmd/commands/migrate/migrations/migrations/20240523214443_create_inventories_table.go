package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateInventoriesTable, downCreateInventoriesTable)
}

func upCreateInventoriesTable(ctx context.Context, tx *sql.Tx) error {
	var queries = `CREATE TABLE IF NOT EXISTS inventories (
        id SERIAL,
        warehouse_id INT REFERENCES warehouses(id),
        product_id INT REFERENCES products(id),
        quantity INT NOT NULL CHECK (quantity > 0),
        reserved_quantity INT NOT NULL DEFAULT 0,
        available_quantity INT GENERATED ALWAYS AS (quantity - reserved_quantity) STORED,
        PRIMARY KEY (warehouse_id, product_id),
        CONSTRAINT valid_reserved_quantity CHECK (quantity >= reserved_quantity)
    );`

	_, err := tx.Exec(queries)
	if err != nil {
		return fmt.Errorf("error exec query. %w", err)
	}

	return nil
}

func downCreateInventoriesTable(ctx context.Context, tx *sql.Tx) error {
	var query = `DROP TABLE IF EXISTS inventories`
	_, err := tx.Exec(query)
	if err != nil {
		return fmt.Errorf("error drop inventories table. %w", err)
	}

	return nil
}
