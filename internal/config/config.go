package config

type Config struct {
	HttpListenaddr int
	IsDev          bool
	DatabaseDsn    string
}

// TODO use env variables
func NewConfig() Config {
	return Config{
		HttpListenaddr: 3000,
		IsDev:          true,
		DatabaseDsn:    "postgres://user:111111@warehouse-postgres:5432/warehousedb",
	}
}
