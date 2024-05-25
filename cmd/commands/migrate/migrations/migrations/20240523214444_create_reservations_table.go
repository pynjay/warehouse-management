package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateReservationsTable, downCreateReservationsTable)
}

func upCreateReservationsTable(ctx context.Context, tx *sql.Tx) error {
    var queries = `CREATE TABLE IF NOT EXISTS reservations (
        reservation_id SERIAL PRIMARY KEY,
        warehouse_id INT REFERENCES warehouses(id),
        product_id INT REFERENCES products(id),
        quantity INT NOT NULL,
        status VARCHAR(20) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	_, err := tx.Exec(queries)
	if err != nil {
		return fmt.Errorf("error exec query. %w", err)
	}

	return nil
}

func downCreateReservationsTable(ctx context.Context, tx *sql.Tx) error {
	var query = `DROP TABLE IF EXISTS reservations`
	_, err := tx.Exec(query)
	if err != nil {
		return fmt.Errorf("error drop reservations table. %w", err)
	}

	return nil
}
