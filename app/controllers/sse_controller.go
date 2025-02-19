package controllers

import (
	"fmt"
	virtualrooms "kaduhod/video-sync/app/virtual_rooms"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)
type SSEController struct {
   virtualRoomsManager *virtualrooms.VirtualRoomsManager
   sseManager *virtualrooms.SSEManager
}
func NewSSEController(manager *virtualrooms.VirtualRoomsManager, sseManager *virtualrooms.SSEManager) *SSEController {
   return &SSEController{
      virtualRoomsManager: manager,
      sseManager: sseManager,
   }
}
func (self *SSEController) getSession(c echo.Context) (*sessions.Session, error) {
    session, err := session.Get("session", c)
    if err != nil {
        return session, err
    }
    return session, nil
}
func (self *SSEController) StreamRoom(c echo.Context) error {
    session, err := self.getSession(c)
    if err != nil {
        fmt.Println(err)
        return c.String(400, err.Error())
    }
    userName := session.Values["user_name"].(string)
    user, err := self.virtualRoomsManager.GetUser(userName)
    if err != nil {
        fmt.Println(err)
        return c.String(400, err.Error())
    }
    roomName := c.Param("roomName")
    fmt.Println(roomName, user)
    room, err := self.virtualRoomsManager.GetRoom(roomName)
    if err != nil {
        fmt.Println(err)
        return c.String(400, err.Error())
    }
    user.Ctx = &c
    if room.AdminId == user.Id {
        self.sseManager.AddRoomSSE(&room)
    }
    self.sseManager.AddUserSSE(&c, &user)
    c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
    c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
    c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

    <-c.Request().Context().Done()
    return nil
}
func (self *SSEController) Post(c echo.Context) error {
    userName := c.FormValue("user_name")
    user, err := self.virtualRoomsManager.GetUser(userName)
    if err != nil {
        fmt.Println(err)
        return c.String(400, err.Error())
    }
    roomName := c.Param("roomName")
    room, err := self.virtualRoomsManager.GetRoom(roomName)
    if err != nil {
        fmt.Println(err)
        return c.String(400, err.Error())
    }
    action := c.FormValue("action")
    fmt.Println("params", action, user, room)
    return c.String(200, "OK")
}
