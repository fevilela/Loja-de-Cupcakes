package controllers

import (
	"github.com/fevilela/cupcakestore/config"
	"github.com/fevilela/cupcakestore/models"
	"github.com/fevilela/cupcakestore/services"
	"github.com/fevilela/cupcakestore/session"
	"github.com/fevilela/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	RenderLogin(ctx *fiber.Ctx) error
	RenderRegister(ctx *fiber.Ctx) error
}

type authController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

func (c *authController) Register(ctx *fiber.Ctx) error {
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return views.RenderError(ctx, "auth/register", nil, "Dados da conta inválidos: "+err.Error())
	}

	profile := &models.Profile{
		FirstName: ctx.FormValue("firstname"),
		LastName:  ctx.FormValue("lastname"),
		User:      *user,
	}

	if err := c.authService.Register(profile); err != nil {
		return views.RenderError(ctx, "auth/register", nil, "Falha ao criar usuário: "+err.Error())
	}

	return ctx.Redirect("/auth/login")
}

func (c *authController) Login(ctx *fiber.Ctx) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	if err := c.authService.Authenticate(ctx, email, password); err != nil {
		return views.RenderError(ctx, "auth/login", nil, "Credenciais inválidas ou usuário inativo.")
	}

	return ctx.Redirect(config.Instance().GetEnvVar("REDIRECT_AFTER_LOGIN", "/"))
}

func (c *authController) Logout(ctx *fiber.Ctx) error {
	sess, err := session.Store.Get(ctx)
	if err != nil {
		return err
	}

	if err := sess.Destroy(); err != nil {
		return err
	}

	return ctx.Redirect(config.Instance().GetEnvVar("REDIRECT_AFTER_LOGOUT", "/"))
}

func (c *authController) RenderLogin(ctx *fiber.Ctx) error {
	return views.Render(ctx, "auth/login", nil)
}

func (c *authController) RenderRegister(ctx *fiber.Ctx) error {
	return views.Render(ctx, "auth/register", nil)
}
