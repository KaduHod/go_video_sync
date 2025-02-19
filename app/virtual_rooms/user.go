package virtualrooms

import (
	"github.com/labstack/echo/v4"
)
type VideoUpdate struct {
    Action string `json:"action"`
    SenderId string `json:"sender_id"`
    RoomId int `json:"room_id"`
}
type User struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Ctx *echo.Context `json:"-"`
}
func (self *User) CloseSSE() error {
    ctx := *self.Ctx
    return ctx.String(200, "Closed")
}
