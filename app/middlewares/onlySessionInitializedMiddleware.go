package middlewares

import (
	"fmt"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func OnlySessionInitialized(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        session, err := session.Get("session", c)
        if err != nil || session == nil {
            url := fmt.Sprintf("/?error=%s", "Unauthorized")
            return c.Redirect(303, url)
        }
        sessionUserId := session.Values["user_id"]
        if sessionUserId == nil {
            url := fmt.Sprintf("/?error=%s", "Unauthorized")
            return c.Redirect(303, url)
        }
        return next(c)
    }
}

