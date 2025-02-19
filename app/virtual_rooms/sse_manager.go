package virtualrooms

import (
	"encoding/json"
	"fmt"

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
    self.Users[user.Id] = user
}
func (self *SSEManager) SSEMessageFromJson(data string) (SSEMessage, error) {
    var message SSEMessage
    err := json.Unmarshal([]byte(data), &message)
    return message, err
}
func (self *SSEManager) AddRoomSSE(room *Room) {
    self.Rooms[room.Id] = room
}
func (self *SSEManager) Send(value string, userDest *User) {
    // identificar quem enviou
    _, err := SSEMessageFromJson(value)
    if err != nil {
        fmt.Println(err)
        return
    }
    context := *userDest.Ctx
    _, err = context.Response().Write([]byte(value))
    if err != nil {
        fmt.Println(err)
    }
}
func (self *SSEManager) StartRoom(room *Room) {
    for {
        select {
        case value, ok := <- room.listener:
            if !ok {
                fmt.Println("Problema com o canal >> ", value)
            } else {
                message, err := SSEMessageFromJson(value)
                if err != nil {
                    fmt.Println(err)
                } else {
                    var destId string
                    if room.AdminId == message.SenderId {
                        destId = room.GuestId
                    } else {
                        destId = room.AdminId
                    }
                    dest := self.Users[destId]
                    admin := self.Users[room.AdminId]
                    self.Send(value, dest)
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
    SenderId string `json:"sender_id"`
    RoomName string `json:"room_name"`
}
func (self *SSEMessage) ToJson() (string, error) {
    json, err := json.Marshal(self)
    if err != nil {
        return "", err
    }
    return string(json), nil
}
