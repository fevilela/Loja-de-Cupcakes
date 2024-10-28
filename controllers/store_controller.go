package controllers

import (
	"github.com/fevilela/cupcakestore/models"
	"github.com/fevilela/cupcakestore/services"
	"github.com/fevilela/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type StoreController interface {
	RenderStore(ctx *fiber.Ctx) error
}

type storeController struct {
	productService services.ProductService
}

func NewStoreController(productService services.ProductService) StoreController {
	return &storeController{
		productService: productService,
	}
}

func (c *storeController) RenderStore(ctx *fiber.Ctx) error {
	query := ctx.Query("q", "")
	page := ctx.QueryInt("page")
	limit := ctx.QueryInt("limit")

	filter := models.NewProductFilter(query, page, limit)
	products := c.productService.FindActiveWithStock(filter)

	data := fiber.Map{
		"Products": products,
		"Filter":   filter,
	}

	return views.Render(ctx, "store/store", data, views.StoreLayout)
}
