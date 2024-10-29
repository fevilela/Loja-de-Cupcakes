package routers

import (
	"github.com/fevilela/cupcakestore/controllers"
	"github.com/fevilela/cupcakestore/database"
	"github.com/fevilela/cupcakestore/repositories"
	"github.com/fevilela/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
)

type StoreRouter struct {
	storeController controllers.StoreController
}

func NewStoreRouter() *StoreRouter {
	// Initialize repositories
	productRepository := repositories.NewProductRepository(database.DB)

	// Initialize services with repositories
	productService := services.NewProductService(productRepository)

	// Initialize controllers with services
	storeController := controllers.NewStoreController(productService)
	return &StoreRouter{
		storeController: storeController,
	}
}

func (r *StoreRouter) InstallRouters(app *fiber.App) {
	store := app.Group("/store")
	store.Get("/", r.storeController.RenderStore)
}
