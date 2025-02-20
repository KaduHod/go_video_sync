package virtualrooms

import (
	"html/template"
)
type Room struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Password string `json:"password"`
    GuestLink template.URL `json:"guest_link"`
    listener chan SSEMessage `json:"-"`
    Guest User `json:"guest_ctx"`
    Admin User `json:"admin_ctx"`
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
    close(self.listener)
}
func (self *Room) SetGuest(guest User) {
    self.Guest = guest
}
