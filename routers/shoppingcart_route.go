package routers

import (
	"github.com/fevilela/cupcakestore/controllers"
	"github.com/fevilela/cupcakestore/database"
	"github.com/fevilela/cupcakestore/middlewares"
	"github.com/fevilela/cupcakestore/repositories"
	"github.com/fevilela/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
)

type ShoppingCartRouter struct {
	shoppingCartController controllers.ShoppingCartController
}

func NewShoppingCartRouter() *ShoppingCartRouter {
	// Initialize repositories
	shoppingCartRepository := repositories.NewShoppingCartRepository(database.DB)
	shoppingCartItemRepository := repositories.NewShoppingCartItemRepository(database.DB)

	// Initialize services with repositories
	shoppingCartItemService := services.NewShoppingCartItemService(shoppingCartItemRepository)
	shoppingCartService := services.NewShoppingCartService(shoppingCartRepository, shoppingCartItemService)

	// Initialize controllers with services
	shoppingCartController := controllers.NewShoppingCartController(shoppingCartService)

	return &ShoppingCartRouter{
		shoppingCartController: shoppingCartController,
	}
}

func (r *ShoppingCartRouter) InstallRouters(app *fiber.App) {
	cart := app.Group("/cart").Use(middlewares.LoginRequired())
	cart.Get("/", r.shoppingCartController.RenderShoppingCart)
	cart.Post("/", r.shoppingCartController.AddShoppingCartItem)
	cart.Get("/count", r.shoppingCartController.CountShoppingCart)
	cart.Get("/remove/:id", r.shoppingCartController.RemoveFromCart)
}
