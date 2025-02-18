package virtualrooms
type VideoUpdate struct {
    Action string `json:"action"`
    SenderId string `json:"sender_id"`
    RoomId int `json:"room_id"`
}
type User struct {
    Id string `json:"id"`
    Name string `json:"name"`
}
