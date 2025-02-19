package virtualrooms

import (
	"html/template"
)
type Room struct {
    Id string `json:"id"`
    Name string `json:"name"`
    AdminId string `json:"admin"`
    GuestId string `json:"guest"`
    Password string `json:"password"`
    GuestLink template.URL `json:"guest_link"`
    listener chan SSEMessage
}
func (self *Room) GetListener() <-chan SSEMessage {
    return self.listener
}
func (self *Room) GetListenerWrite() chan <- SSEMessage {
    return self.listener
}
func (self *Room) SetListener(listener chan SSEMessage) {
    self.listener = listener
}
