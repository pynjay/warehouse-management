//go:build wireinject
// +build wireinject

package factory

import (
	"context"

	"fmt"
	"warehouse/internal/infrastructure/http"

	"github.com/google/wire"
	"gorm.io/gorm"
)

type Service struct {
}

func InitializeService(ctx context.Context) (Service, func(), error) {
	panic(wire.Build(
		handlerSet,
		provideConfig,
		provideService,
		infrustructureSet,
		interfacesSet,
		usecaseSet,
		repositorySet,
	))
}

func InitializeMigrationContainer() (MigrationContainer, func(), error) {
	panic(wire.Build(
		provideConfig,
		provideGorm,
		provideMigrationContainer,
	))
}

func provideMigrationContainer(gormDb *gorm.DB) (MigrationContainer, error) {
	db, err := gormDb.DB()
	if err != nil {
		return MigrationContainer{}, fmt.Errorf("error get db from gorm. %w", err)
	}
	DefaultMigrationContainer = MigrationContainer{db: db}

	return DefaultMigrationContainer, nil
}

func provideService(
	ctx context.Context,
	_ *http.HTTPServer,
) Service {
	return Service{}
}
