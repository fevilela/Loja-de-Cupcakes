package routers

import (
	"github.com/fevilela/cupcakestore/controllers"
	"github.com/fevilela/cupcakestore/database"
	"github.com/fevilela/cupcakestore/middlewares"
	"github.com/fevilela/cupcakestore/repositories"
	"github.com/fevilela/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
)

type AuthRouter struct {
	authController controllers.AuthController
}

func NewAuthRouter() *AuthRouter {
	// Initialize repositories
	userRepository := repositories.NewUserRepository(database.DB)
	profileRepository := repositories.NewProfileRepository(database.DB)

	// Initialize services with repositories
	userService := services.NewUserService(userRepository)
	profileService := services.NewProfileService(profileRepository)
	authService := services.NewAuthService(userService, profileService)

	// Initialize controllers with services
	authController := controllers.NewAuthController(authService)

	return &AuthRouter{
		authController: authController,
	}
}

func (r *AuthRouter) InstallRouters(app *fiber.App) {
	auth := app.Group("/auth")

	auth.Get("/login", r.authController.RenderLogin)
	auth.Post("/login", r.authController.Login)
	auth.Get("/register", r.authController.RenderRegister)
	auth.Post("/register", r.authController.Register)
	auth.Get("/logout", r.authController.Logout).Use(middlewares.LoginRequired())
}
