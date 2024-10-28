package controllers

import (
	"github.com/fevilela/cupcakestore/models"
	"github.com/fevilela/cupcakestore/services"
	"github.com/fevilela/cupcakestore/utils"
	"github.com/fevilela/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	RenderCreate(ctx *fiber.Ctx) error
	RenderUsers(ctx *fiber.Ctx) error
	RenderUser(ctx *fiber.Ctx) error
	RenderDelete(ctx *fiber.Ctx) error
}

type userController struct {
	userService services.UserService
}

func NewUserController(u services.UserService) UserController {
	return &userController{
		userService: u,
	}
}

func (c *userController) RenderCreate(ctx *fiber.Ctx) error {
	return views.Render(ctx, "users/create", nil, views.BaseLayout)
}

func (c *userController) Create(ctx *fiber.Ctx) error {
	user := &models.User{}
	if err := ctx.BodyParser(user); err != nil {
		return views.RenderError(ctx, "users/create", nil, "Dados de usuário inválidos: "+err.Error(), views.BaseLayout)
	}

	user.IsStaff = ctx.FormValue("isStaff") == "on"
	user.IsActive = ctx.FormValue("isActive") == "on"

	if err := c.userService.Create(user); err != nil {
		return views.RenderError(ctx, "users/create", nil, "Falha ao criar usuário: "+err.Error(), views.BaseLayout)
	}

	return ctx.Redirect("/users")
}

func (c *userController) RenderUsers(ctx *fiber.Ctx) error {
	query := ctx.Query("q", "")
	page := ctx.QueryInt("page")
	limit := ctx.QueryInt("limit")
	filter := models.NewUserFilter(query, page, limit)
	users := c.userService.FindAll(filter)

	return views.Render(ctx, "users/users", fiber.Map{"Users": users, "Filter": filter}, views.BaseLayout)
}

func (c *userController) RenderUser(ctx *fiber.Ctx) error {
	user, err := c.getUser(ctx)
	if err != nil {
		return err
	}

	userSess, err := c.getUserSession(ctx)
	if err != nil {
		return err
	}

	if user.IsStaff && user.ID != userSess.ID {
		return ctx.Redirect("/users")
	}

	layout := selectLayout(userSess.IsStaff, user.ID == userSess.ID)
	if layout == "" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	return views.Render(ctx, "users/user", user, layout)
}

func (c *userController) getUser(ctx *fiber.Ctx) (models.User, error) {
	userID, err := utils.StringToId(ctx.Params("id"))
	if err != nil {
		return models.User{}, ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return c.userService.FindById(userID)
}

func (c *userController) getUserSession(ctx *fiber.Ctx) (*models.User, error) {
	userSess, ok := ctx.Locals("Profile").(*models.Profile)
	if !ok || userSess == nil {
		return nil, fiber.ErrUnauthorized
	}
	return &userSess.User, nil
}

func (c *userController) Update(ctx *fiber.Ctx) error {
	user, err := c.getUserAndCheckAccess(ctx)
	if err != nil {
		return err
	}

	if err := c.updateUserFromRequest(ctx, user); err != nil {
		return views.RenderError(ctx, "users/user", user, err.Error(), selectLayout(user.IsStaff, user.ID == ctx.Locals("Profile").(*models.Profile).UserID))
	}

	if err := c.updateUserPassword(ctx, user); err != nil {
		return views.RenderError(ctx, "users/user", user,
			"Falha ao atualizar a senha. Certifique-se de que está inserido corretamente.", selectLayout(user.IsStaff, user.ID == ctx.Locals("Profile").(*models.Profile).UserID))
	}

	if err := c.userService.Update(user); err != nil {
		return views.RenderError(ctx, "users/user", user,
			"Falha ao atualizar o usuário.", selectLayout(user.IsStaff, user.ID == ctx.Locals("Profile").(*models.Profile).UserID))
	}

	if user.ID == ctx.Locals("Profile").(*models.Profile).UserID {
		return ctx.Redirect("/auth/logout")
	}

	redirectPath := selectRedirectPath(user.IsStaff)
	return ctx.Redirect(redirectPath)
}

func (c *userController) getUserAndCheckAccess(ctx *fiber.Ctx) (*models.User, error) {
	id, err := utils.StringToId(ctx.Params("id"))
	if err != nil {
		return nil, ctx.Redirect("/users")
	}

	user, err := c.userService.FindById(id)
	if err != nil {
		return nil, ctx.Redirect("/users")
	}

	userSess, err := c.getUserSession(ctx)
	if err != nil {
		return nil, err
	}

	if user.IsStaff && user.ID != userSess.ID {
		return nil, ctx.SendStatus(fiber.StatusUnauthorized)
	}

	if !userSess.IsStaff && user.ID != userSess.ID {
		return nil, ctx.SendStatus(fiber.StatusUnauthorized)
	}

	return &user, nil
}

func (c *userController) updateUserFromRequest(ctx *fiber.Ctx, user *models.User) error {
	if err := ctx.BodyParser(user); err != nil {
		return err
	}

	user.IsStaff = ctx.FormValue("isStaff") == "on"
	user.IsActive = ctx.FormValue("isActive") == "on"

	return nil
}

func (c *userController) updateUserPassword(ctx *fiber.Ctx, user *models.User) error {
	oldPassword := ctx.FormValue("oldPassword")
	newPassword := ctx.FormValue("newPassword")

	if oldPassword != "" && newPassword != "" {
		err := user.UpdatePassword(oldPassword, newPassword)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *userController) RenderDelete(ctx *fiber.Ctx) error {
	user, err := c.getUser(ctx)
	if err != nil {
		return ctx.Redirect("/users")
	}

	return views.Render(ctx, "users/delete", user, views.BaseLayout)
}

func (c *userController) Delete(ctx *fiber.Ctx) error {
	id, err := utils.StringToId(ctx.Params("id"))
	if err != nil {
		return ctx.Redirect("/users")
	}

	err = c.userService.Delete(id)
	if err != nil {
		return ctx.Redirect("/users")
	}

	return ctx.Redirect("/users")
}
