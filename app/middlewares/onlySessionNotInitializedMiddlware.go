package middlewares

import (
	"regexp"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func OnlySessionNotInitialized(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        session, err := session.Get("session", c)
        if err == nil && session.Values["user_id"] != nil {
            room, linkDeConvidado := extractRoomName(c.Request().RequestURI)
            if linkDeConvidado && c.Request().Method == "GET" {
                return c.Redirect(303, "/app/join/room/"+room)
            }
        }
        return next(c)
    }
}
func extractRoomName(path string) (string, bool) {
	re := regexp.MustCompile(`^/join/room/([^/]+)$`)
	matches := re.FindStringSubmatch(path)

	if len(matches) > 1 {
		return matches[1], true
	}
	return "", false
}
