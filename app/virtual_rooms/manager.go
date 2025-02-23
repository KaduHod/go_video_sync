package virtualrooms

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"kaduhod/video-sync/app/db"
	"kaduhod/video-sync/app/utils"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)
type VirtualRoomsManager struct {
    Rooms []*Room
    Users []*User
    UsersLimit int
    RoomsLimit int
    DefaultExpirationTime time.Duration
    RedisConn *redis.Client
    Ctx context.Context
}
func NewVirtualRoomsManager() *VirtualRoomsManager {
    ctx := context.Background()
    redisConn, err := db.RedisConn(ctx)
    if err != nil {
        panic(err)
    }
    var users []*User
    var rooms []*Room
    return &VirtualRoomsManager{
        Rooms: rooms,
        Users: users,       UsersLimit: 10,
        RoomsLimit: 5,
        DefaultExpirationTime: time.Hour * 2,
        RedisConn: redisConn,
        Ctx: ctx,
    }
}
func (self *VirtualRoomsManager) Init() {
    var userAdmin User
    exists, err := self.UserExists("Admin")
    if err != nil {
        panic(err)
    }
    if !exists {
        userAdmin := User {
            Id: uuid.New().String(),
            Name: "Admin",
        }
        fmt.Println("Criando admin")
        if err = self.AddUser(userAdmin); err != nil {
            panic(err)
        }
    } else {
        userAdmin, err = self.GetUser("Admin")
        if err != nil {
            panic(err)
        }
    }
    roomExists, err := self.RoomExists("Teste")
    if err != nil {
        panic(err)
    }
    if !roomExists {
        pass, err := utils.HashPassword("123456")
        if err != nil {
            panic(err)
        }
        testRoom := Room {
            Id: uuid.New().String(),
            Name: "Teste",
            Password: pass,
            GuestLink: template.URL("http://localhost:3003/join/room/Teste"),
            Admin: &userAdmin,
        }
        if err := self.AddRoom(testRoom); err != nil {
            panic(err)
        }
    }
}
func (self *VirtualRoomsManager) AddRoom(room Room) error {
    roomExists, err := self.RoomExists(room.Name)
    if err != nil {
        return err
    }
    if roomExists {
        return errors.New("Room key already exists >> room:" + room.Name)
    }
    roomStr := utils.JsonEncode(room)
    _, err = self.RedisConn.Set(self.Ctx ,"room:"+room.Name, roomStr, self.DefaultExpirationTime).Result()
    if err != nil {
        return err
    }
    return nil
}
func (self *VirtualRoomsManager) AddUser(user User) error {
    userExists, err := self.UserExists(user.Name)
    if err != nil {
        return err
    }
    if userExists {
        return errors.New("User key already exists >> user:" + user.Name)
    }
    userStr := utils.JsonEncode(user)
    _, err = self.RedisConn.Set(self.Ctx ,"user:"+user.Name, userStr, self.DefaultExpirationTime).Result()
    if err != nil {
        return err
    }
    return nil
}
func (self *VirtualRoomsManager) GuestJoinRoom(guest User, password string, roomName string) error {
    userExists, err := self.UserExists(guest.Name)
    if err != nil {
        return err
    }
    if !userExists {
        return errors.New("User does not exist")
    }
    roomExists, err := self.RoomExists(roomName)
    if err != nil {
        return err
    }
    if !roomExists {
        return errors.New("Room does not exist")
    }
    roomRedis, err := self.GetRoom(roomName)
    if err != nil {
        return err
    }
    if roomRedis.Guest != nil {
        return errors.New("Room is already full")
    }
    if !utils.CheckPasswordHash(password, roomRedis.Password) {
        return errors.New("Wrong password")
    }
    roomRedis.Guest = &guest
    if err := self.RedisConn.Set(self.Ctx, "room:"+roomName, utils.JsonEncode(roomRedis), self.DefaultExpirationTime).Err(); err != nil {
        return err
    }
    return nil
}
func (self *VirtualRoomsManager) RoomExists(roomName string) (bool, error) {
    exists, err := self.RedisConn.Exists(self.Ctx, "room:"+roomName).Result()
    if err != nil {
        return false, err
    }
    if exists == 1 {
        return true, nil
    }
    return false, nil
}
func (self *VirtualRoomsManager) UserExists(userName string) (bool, error) {
    exists, err := self.RedisConn.Exists(self.Ctx, "user:"+userName).Result()
    if err != nil {
        return false, err
    }
    if exists == 1 {
        return true, nil
    }
    return false, nil
}
func (self *VirtualRoomsManager) GetRoom(roomName string) (Room, error) {
    var room Room
    roomExists, err := self.RoomExists(roomName)
    if err != nil {
        return room, err
    }
    if !roomExists {
        return room, errors.New("Room does not exist")
    }
    roomStr, err := self.RedisConn.Get(self.Ctx, "room:"+roomName).Result()
    if err != nil {
        return room, err
    }
    utils.JsonDecode(roomStr, &room)
    return room, nil
}
func (self *VirtualRoomsManager) DeleteKey(key string) {
   self.RedisConn.Del(self.Ctx, key)
}
func (self *VirtualRoomsManager) GetUser(userName string) (User, error) {
    var user User
    exists, err := self.UserExists(userName)
    if err != nil {
        return user, err
    }
    if !exists {
        return user, errors.New("User does not exist")
    }
    userStr, err := self.RedisConn.Get(self.Ctx, "user:"+userName).Result()
    if err != nil {
        return user, err
    }
    utils.JsonDecode(userStr, &user)
    return user, nil
}
func (self *VirtualRoomsManager) RemoveUserFromRoom(user User, room Room) error {
    roomExists, err := self.RoomExists(room.Name)
    if err != nil {
        return err
    }
    if !roomExists {
        return errors.New("Room does not exist")
    }
    userExists, err := self.UserExists(user.Name)
    if err != nil {
        return err
    }
    if !userExists {
        return errors.New("User does not exist")
    }
    room, err = self.GetRoom(room.Name)
    if err != nil {
        return err
    }
    user, err = self.GetUser(user.Name)
    if err != nil {
        return err
    }
    if user.Id == room.Admin.Id {
        room.Admin = nil
    } else {
        room.Guest = nil
    }
    roomStr := utils.JsonEncode(room)
    _, err = self.RedisConn.Set(self.Ctx ,"room:"+room.Name, roomStr, self.DefaultExpirationTime).Result()
    if err != nil {
        return err
    }
    return nil
}
