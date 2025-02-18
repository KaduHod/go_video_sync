package tests

import (
	"fmt"
	virtualrooms "kaduhod/video-sync/app/virtual_rooms"
	"testing"
)

func TestMain(t *testing.T) {
    m := virtualrooms.NewVirtualRoomsManager()
    m.Init()
}
func TestGet(t *testing.T) {
    m := virtualrooms.NewVirtualRoomsManager()
    m.Init()
    _, err := m.GetRoom("Test")
    if err != nil {
        t.Error(err)
    }
    _, err = m.GetUser("Admin")
    if err != nil {
        t.Error(err)
    }
}
func TestMissGet(t *testing.T) {
    m := virtualrooms.NewVirtualRoomsManager()
    m.Init()
    _, err := m.GetUser("inexistente")
    if err == nil {
        t.Fail()
    }
}
func TestRedisKeyExists(t *testing.T) {
    m := virtualrooms.NewVirtualRoomsManager()
    result, err := m.RedisConn.Exists(m.Ctx, "room:Test").Result()
    fmt.Println(result, err, "Aqui")
}
