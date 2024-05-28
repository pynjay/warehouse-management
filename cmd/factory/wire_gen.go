// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package factory

import (
	"context"
	"fmt"
	"warehouse/internal/handlers/apiv1/products"
	"warehouse/internal/handlers/apiv1/reservations"
	"warehouse/internal/handlers/apiv1/warehouses"
	"warehouse/internal/infrastructure/http"
	"warehouse/internal/infrastructure/repository"
	"warehouse/internal/presenters"
	"warehouse/internal/usecases"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func InitializeService(ctx context.Context) (Service, func(), error) {
	config := provideConfig()
	logger, err := provideLogger(config)
	if err != nil {
		return Service{}, nil, err
	}
	validate := validator.New()
	db, err := provideGorm(config)
	if err != nil {
		return Service{}, nil, err
	}
	inventoriesRepositoryPostgres := repository.NewInventoriesRepositoryPostgres(db)
	reservationsRepositoryPostgres := repository.NewReservationsRepositoryPostgres(db)
	productsRepositoryPostgres := repository.NewProductsRepositoryPostgres(db, config)
	makeReservationsUseCaseImpl := usecases.NewMakeReservationsUseCase(inventoriesRepositoryPostgres, reservationsRepositoryPostgres, productsRepositoryPostgres, db)
	reservationsApiPresenterImpl := presenters.NewReservationsApiPresenterImpl()
	updateReservationStatusUseCaseImpl := usecases.NewUpdateReservationStatusUseCase(db, reservationsRepositoryPostgres, inventoriesRepositoryPostgres)
	handler := reservations.NewHandler(logger, validate, makeReservationsUseCaseImpl, reservationsApiPresenterImpl, updateReservationStatusUseCaseImpl)
	createProductsUseCaseImpl := usecases.NewCreateProductsUseCase(productsRepositoryPostgres)
	productsApiPresenterImpl := presenters.NewProductsApiPresenterImpl()
	productsHandler := products.NewHandler(logger, validate, createProductsUseCaseImpl, productsApiPresenterImpl, inventoriesRepositoryPostgres)
	warehousesRepositoryPostgres := repository.NewWarehousesRepositoryPostgres(db)
	createWarehousesUseCaseImpl := usecases.NewCreateWarehousesUseCase(warehousesRepositoryPostgres)
	warehousesHandler := warehouses.NewHandler(logger, validate, createWarehousesUseCaseImpl, inventoriesRepositoryPostgres)
	httpServer, cleanup := provideHTTPServer(logger, config, handler, productsHandler, warehousesHandler)
	service := provideService(ctx, httpServer)
	return service, func() {
		cleanup()
	}, nil
}

func InitializeMigrationContainer() (MigrationContainer, func(), error) {
	config := provideConfig()
	db, err := provideGorm(config)
	if err != nil {
		return MigrationContainer{}, nil, err
	}
	migrationContainer, err := provideMigrationContainer(db)
	if err != nil {
		return MigrationContainer{}, nil, err
	}
	return migrationContainer, func() {
	}, nil
}

// wire.go:

type Service struct {
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