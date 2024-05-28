package factory

import (
	"warehouse/internal/config"
	repository2 "warehouse/internal/infrastructure/repository"
	"warehouse/internal/repository"

	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var repositorySet = wire.NewSet(
	provideGorm,
	repository2.NewProductsRepositoryPostgres,
	wire.Bind(new(repository.ProductsRepository), new(*repository2.ProductsRepositoryPostgres)),
	repository2.NewWarehousesRepositoryPostgres,
	wire.Bind(new(repository.WarehousesRepository), new(*repository2.WarehousesRepositoryPostgres)),
	repository2.NewInventoriesRepositoryPostgres,
	wire.Bind(new(repository.InventoriesRepository), new(*repository2.InventoriesRepositoryPostgres)),
	repository2.NewReservationsRepositoryPostgres,
	wire.Bind(new(repository.ReservationsRepository), new(*repository2.ReservationsRepositoryPostgres)),
)

func provideGorm(c config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(c.DatabaseDsn), &gorm.Config{})
}
