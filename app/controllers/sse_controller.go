package controllers

import (
	"fmt"
	virtualrooms "kaduhod/video-sync/app/virtual_rooms"
	"time"

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
    room, err := self.virtualRoomsManager.GetRoom(roomName)
    if err != nil {
        fmt.Println(err)
        return c.String(400, err.Error())
    }
    user.Ctx = &c
    if room.AdminId == user.Id {
        listener := make(chan virtualrooms.SSEMessage)
        room.SetListener(listener)
        self.sseManager.AddRoomSSE(&room)
        go self.sseManager.StartRoom(&room)
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
    // sender
    user, err := self.virtualRoomsManager.GetUser(userName)
    if err != nil {
        fmt.Println(err)
        return c.String(400, err.Error())
    }
    roomName := c.Param("roomName")
    // roomDest
    room, ok := self.sseManager.Rooms[roomName]
    if room == nil {
        fmt.Println("Room not found")
        return c.String(400, "Room not found")
    }
    if !ok {
        fmt.Println(err)
        return c.String(400, err.Error())
    }
    if room.GetListener() == nil {
        fmt.Println("Channel[listener] not initialized")
        return c.String(400, "Channel not initialized")
    }
    msg := virtualrooms.SSEMessage{}
    msg.RoomName = roomName
    msg.Value = c.FormValue("action")
    msg.Sender = user
    msg.Timestamp = time.Now()
    room.GetListenerWrite() <- msg
    return c.String(200, "OK")
}
