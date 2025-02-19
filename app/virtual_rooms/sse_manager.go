package virtualrooms

import "github.com/labstack/echo/v4"

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
func (self *SSEManager) AddRoomSSE(room *Room) {
    room.StartChannels()
    self.Rooms[room.Name] = room
}
