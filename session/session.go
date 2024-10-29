package session

import (
	"time"

	"github.com/fevilela/cupcakestore/models"

	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store *session.Store

func SetupSession() {
	sessConfig := session.Config{
		Expiration: 1 * time.Hour,
	}
	Store = session.New(sessConfig)
	Store.RegisterType(&models.Profile{})
}
