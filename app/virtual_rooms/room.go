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
    listener <-chan string
}
