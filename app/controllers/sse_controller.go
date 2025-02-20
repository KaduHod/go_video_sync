package controllers

import (
	"fmt"
	virtualrooms "kaduhod/video-sync/app/virtual_rooms"
	"time"

	"github.com/labstack/echo/v4"
)
type SSEController struct {
    Controller
    virtualRoomsManager *virtualrooms.VirtualRoomsManager
    sseManager *virtualrooms.SSEManager
}
func NewSSEController(manager *virtualrooms.VirtualRoomsManager, sseManager *virtualrooms.SSEManager) *SSEController {
   return &SSEController{
      virtualRoomsManager: manager,
      sseManager: sseManager,
   }
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
    if room.Admin.Id == user.Id {
        listener := make(chan virtualrooms.SSEMessage)
        room.SetListener(listener)
        self.sseManager.AddRoomSSE(&room)
        go self.sseManager.StartRoom(&room)
    } else {
        existentRoom := self.sseManager.Rooms[roomName]
        if existentRoom == nil {
            fmt.Println("Room not found")
            return c.String(400, "Room not found")
        }
        existentRoom.SetGuest(user)
    }
    self.sseManager.AddUserSSE(&c, &user)
    fmt.Println("Adiciondo User: ", user.Name, " na room: ", room.Name)
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
    room.GetListenerWriter() <- msg
    return c.String(200, "OK")
}
func (self *SSEController) CloseRoom(c echo.Context) error {
    roomRedis, err := self.virtualRoomsManager.GetRoom(c.Param("roomName"))
    if err != nil {
        fmt.Println(err)
        return c.String(400, err.Error())
    }
    adminUserRedis, err := self.virtualRoomsManager.GetUser(roomRedis.Admin.Name)
    if err != nil {
        fmt.Println(err)
        return c.String(400, err.Error())
    }
    guestUserRedis, err := self.virtualRoomsManager.GetUser(roomRedis.Guest.Name)
    if err != nil {
        fmt.Println(err)
        return c.String(400, err.Error())
    }
    roomRedisSse := self.sseManager.Rooms[roomRedis.Name]
    if roomRedisSse == nil {
        fmt.Println("Room not found")
        return c.String(400, "Room not found")
    }
    guest := self.sseManager.Users[guestUserRedis.Name]
    admin := self.sseManager.Users[adminUserRedis.Name]
    if admin == nil {
        fmt.Println("Admin not found")
        return c.String(400, "Admin not found")
    }
    if guest == nil {
       fmt.Println("Guest not found")
       return c.String(400, "User not found")
    }
    if err := guest.CloseSSE(); err != nil {
        fmt.Println(err)
        return c.String(400, err.Error())
    }
    if err := admin.CloseSSE(); err != nil {
        fmt.Println(err)
        return c.String(400, err.Error())
    }
    room := self.sseManager.Rooms[roomRedis.Name]
    if room == nil {
        fmt.Println("Room not found")
        return c.String(400, "Room not found")
    }
    self.virtualRoomsManager.DeleteKey("roomRedis:"+roomRedis.Name)
    room.Close()
    self.sseManager.DeleteRoom(roomRedis.Name)
    fmt.Println("Room closed >> ", roomRedis.Name)
    return c.String(200, "OK")
}
