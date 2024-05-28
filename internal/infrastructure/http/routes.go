package http

import "warehouse/pkg/middlewares"

func (server *HTTPServer) BuildRouter() {
	server.e.GET("ping", server.handlePing)

	apiV1 := server.e.Group("api/v1/")
	apiV1.Use(middlewares.AcceptJsonMiddleware())
	// TODO auth

	{
		apiV1Products := apiV1.Group("products")

		apiV1Products.POST("", server.productsHandler.HandleCreateProductsTypes)
		apiV1Products.POST("/:product_id", server.productsHandler.HandleAddProducts)
	}

	{
		apiV1Warehouses := apiV1.Group("warehouses")

		apiV1Warehouses.POST("", server.warehousesHandler.HandleCreateWarehouse)
		apiV1Warehouses.GET("/:warehouse_id/count", server.warehousesHandler.HandleTotalQuantityCount)
	}

	{
		apiV1Reservations := apiV1.Group("reservations")

		apiV1Reservations.POST("", server.reservationsHandler.HandleMakeReservations)
		apiV1Reservations.PATCH("/:reservation_id", server.reservationsHandler.HandleUpdateReservationStatus)
	}
}
