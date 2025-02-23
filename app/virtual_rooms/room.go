package virtualrooms

import (
	"fmt"
	"html/template"
	"time"
)
type Room struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Password string `json:"password"`
    GuestLink template.URL `json:"guest_link"`
    listener chan SSEMessage `json:"-"`
    Guest *User `json:"guest_ctx"`
    Admin *User `json:"admin_ctx"`
}
func (self *Room) GetListener() <-chan SSEMessage {
    return self.listener
}
func (self *Room) GetListenerWriter() chan <- SSEMessage {
    return self.listener
}
func (self *Room) SetListener(listener chan SSEMessage) {
    self.listener = listener
}
func (self *Room) Close() {
    fmt.Println("FEchando o canal")
    close(self.listener)
}
func (self *Room) SetGuest(guest *User) {
    self.Guest = guest
}
func (self *Room) SetAdmin(admin *User) {
    self.Admin = admin
}
func (self *Room) AfterUserLeave(user *User) {
    event := SSEMessage{
        Value: "notification::user-leave",
        RoomName: self.Name,
        Sender: *user,
        Timestamp: time.Now(),
    }
    self.GetListenerWriter() <- event
}
func (self *Room) AfterUserJoin(user *User) {
    event := SSEMessage{
        Value: "notification::user-join",
        RoomName: self.Name,
        Sender: *user,
        Timestamp: time.Now(),
    }
    self.GetListenerWriter() <- event
}
