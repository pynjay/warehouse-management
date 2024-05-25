package factory

import "database/sql"

type MigrationContainer struct {
	db *sql.DB
}

func (m MigrationContainer) DB() *sql.DB {
	return m.db
}

var DefaultMigrationContainer MigrationContainer
