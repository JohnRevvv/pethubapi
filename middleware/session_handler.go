package middleware

import (
	"github.com/gofiber/fiber/v2/middleware/session"
)

var SessionStore *session.Store = session.New()
