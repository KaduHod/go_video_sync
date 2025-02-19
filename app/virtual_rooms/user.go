package virtualrooms

import (
	"fmt"

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
    Ctx *echo.Context
}
func (self *User) Send(value string) error {
    ctx := *self.Ctx
    _, err := ctx.Response().Write([]byte(value))
    if err != nil {
        fmt.Println(err)
    }
    return err
}
