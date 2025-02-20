package virtualrooms

import (
	"encoding/json"
	"fmt"
	"kaduhod/video-sync/app/utils"
	"time"

	"github.com/labstack/echo/v4"
)
type SSEManager struct {
    Rooms map[string]*Room
    Users map[string]*User
}
func NewSSEManager() *SSEManager {
    return &SSEManager{
        Rooms: make(map[string]*Room),
        Users: make(map[string]*User),
    }
}
func (self *SSEManager) AddUserSSE(c *echo.Context, user *User) {
    user.Ctx = c
    self.Users[user.Name] = user
}
func (self *SSEManager) DeleteRoom(roomName string) {
    delete(self.Rooms, roomName)
}
func (self *SSEManager) AddRoomSSE(room *Room) {
    self.Rooms[room.Name] = room
}
func (self *SSEManager) Send(value SSEMessage, userDest *User) {
    json, err := utils.JsonStringify(value)
    if err != nil {
        fmt.Println(err)
        return
    }
    context := *userDest.Ctx
    _, err = fmt.Fprintf(context.Response(), "data: %s\n\n", json)
    if err != nil {
        fmt.Println(err)
        return
    }
    context.Response().Flush()
}
func (self *SSEManager) StartRoom(room *Room) {
    for {
        select {
        case value, ok := <- room.GetListener():
            if !ok {
                fmt.Println("Problema com o canal >> ", value)
                return
            } else {
                guest := self.Users[room.Guest.Name]
                if guest != nil {
                    fmt.Println("Enviou para o guest")
                    self.Send(value, guest)
                }
                admin := self.Users[room.Admin.Name]
                if admin != nil {
                    fmt.Println("Enviou para o admin")
                    self.Send(value, admin)
                }
            }
        }
    }
}
func SSEMessageFromJson(data string) (SSEMessage, error) {
    var message SSEMessage
    err := json.Unmarshal([]byte(data), &message)
    return message, err
}
type SSEMessage struct {
    Value string `json:"value"`
    RoomName string `json:"room_name"`
    Sender User `json:"sender"`
    Timestamp time.Time `json:"timestamp"`
}
