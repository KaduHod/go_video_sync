package virtualrooms

import (
	"fmt"
	"kaduhod/video-sync/app/utils"

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
func (self *User) NotifyCloseRoom() {
    event := SSEMessage{
        Value: "close",
    }
    string, err := utils.JsonStringify(event)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    ctx := *self.Ctx
    _, err = fmt.Fprintf(ctx.Response(), "data: %s\n\n", string)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    ctx.Response().Flush()
}
