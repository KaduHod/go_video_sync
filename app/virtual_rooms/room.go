package virtualrooms

import (
	"encoding/json"
	"fmt"
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
    sender chan<- string
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
func SSEMessageFromJson(data string) (SSEMessage, error) {
    var message SSEMessage
    err := json.Unmarshal([]byte(data), &message)
    return message, err
}
func (self *Room) StartChannels() {
    if self.listener != nil || self.sender != nil{
        return
    }
    self.listener = make(<-chan string)
    self.sender = make(chan<- string)
    go self.Listen()
}
func (self *Room) Listen() {
    for {
        select {
        case value, ok := <- self.listener:
            if !ok {
                fmt.Println("Problema com o canal >> ", value)
            } else {
                self.Send(value)
            }
        }
    }
}
func (self *Room) Send(value string, userDest *User) {
    // identificar quem enviou
    msg, err := SSEMessageFromJson(value)
    if err != nil {
        fmt.Println(err)
        return
    }
    context := *userDest.Ctx
    _, err = context.Response().Write([]byte(value))
    if err != nil {
        fmt.Println(err)
    }
    // pegar sse do destintario

    // encaminhar para o destinatario
}
func (self *Room) CloseChannels() {
    close(self.sender)
}
