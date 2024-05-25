package factory

import "warehouse/internal/config"

func provideConfig() config.Config {
	return config.NewConfig()
}

