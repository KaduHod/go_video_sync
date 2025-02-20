package controllers

import (
	"fmt"
	"html/template"
	"kaduhod/video-sync/app/utils"
	virtualrooms "kaduhod/video-sync/app/virtual_rooms"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type VirtualRoomController struct {
    Controller
    virtualRoomsManager *virtualrooms.VirtualRoomsManager
}

func NewVirtualRoomController(manager *virtualrooms.VirtualRoomsManager) *VirtualRoomController {
    return &VirtualRoomController{
        virtualRoomsManager: manager,
    }
}
func (self *VirtualRoomController) getSession(c echo.Context) (*sessions.Session, error) {
    session, err := session.Get("session", c)
    if err != nil {
        return session, err
    }
    return session, nil
}
func (self *VirtualRoomController) defaultErrorReturn(err error, c echo.Context) error {
    fmt.Println(err)
    return c.Redirect(303, "/?error=We are facing some issues, contact the admin")
}
func (self *VirtualRoomController) Index(c echo.Context) error {
    pageData := self.getPageData(c)
    if c.QueryParam("error") != "" {
        pageData["user_error"] = c.QueryParam("error")
    }
    return c.Render(200, "main.tmpl", pageData)
}
func (self *VirtualRoomController) IndexUsers(c echo.Context) error {
    pageData := self.getPageData(c)
    error := c.QueryParam("error")
    if error != "" {
        pageData["user_error"] = error
        return c.Render(400, "user.tmpl", pageData)
    }
    return c.Render(200, "user.tmpl", pageData)
}
func (self *VirtualRoomController) NewUser(c echo.Context) error {
    if len(self.virtualRoomsManager.Users) > self.virtualRoomsManager.UsersLimit {
        return c.Redirect(303, "/?error=Too many users")
    }
    name := c.FormValue("name")
    exists, err := self.virtualRoomsManager.UserExists(name)
    if err != nil {
        return self.defaultErrorReturn(err, c)
    }
    if exists {
        return c.Redirect(303, "/?error=Username not available")
    }

    id := uuid.New().String()
    user := virtualrooms.User {
        Name: name,
        Id: id,
    }
    err = self.virtualRoomsManager.AddUser(user)
    if err != nil {
        return self.defaultErrorReturn(err, c)
    }
    sess, err := self.getSession(c)
    if err != nil {
        return c.Redirect(303, "/?error=" + err.Error())
    }
    sess.Options = &sessions.Options{
        Path:     "/",
        MaxAge:   3600 * 2, // duas horas
        HttpOnly: true,
    }
    sess.Values["user_name"] = user.Name
    sess.Values["user_id"] = user.Id
    if err := sessions.Save(c.Request(), c.Response()); err != nil {
        return c.Redirect(303, "/?error=" + err.Error())
    }
    return c.Redirect(303, "/app/user")
}
func (self *VirtualRoomController) NewRoom(c echo.Context) error {
    pageData := self.getPageData(c)
    if len(self.virtualRoomsManager.Rooms) > self.virtualRoomsManager.RoomsLimit {
        pageData["user_error"] = "Too many rooms"
        return c.Render(400, "main.tmpl", pageData)
    }
    name := c.FormValue("name")
    roomExists, err := self.virtualRoomsManager.RoomExists(name)
    if err != nil {
        return self.defaultErrorReturn(err, c)
    }
    if roomExists {
        return c.Redirect(303, "/app/user?error=Room name unavailable")
    }
    session, err := self.getSession(c)
    if err != nil {
        return self.defaultErrorReturn(err, c)
    }
    userName := session.Values["user_name"].(string)
    userExists, err := self.virtualRoomsManager.UserExists(userName)
    if err != nil || !userExists {
        return self.defaultErrorReturn(err, c)
    }
    password := c.FormValue("password")
    senhaHash, err := utils.HashPassword(password)
    if err != nil {
        return c.Redirect(303, "/app/user?error=" + err.Error())
    }
    user, err := self.virtualRoomsManager.GetUser(userName)
    if err != nil {
        return self.defaultErrorReturn(err, c)
    }
    room := virtualrooms.Room {
        Id: uuid.New().String(),
        Name: name,
        Password: senhaHash,
        GuestLink: template.URL(c.Request().Host+"/join/room/" + name),
        Admin: user,
    }
    if err := self.virtualRoomsManager.AddRoom(room); err != nil {
        return self.defaultErrorReturn(err, c)
    }
    roomUrl := fmt.Sprintf("/app/room/%s", name)
    return c.Redirect(http.StatusSeeOther, roomUrl)
}
func (self *VirtualRoomController) RoomIndex(c echo.Context) error {
    pageData := self.getPageData(c)
    session, err := self.getSession(c)
    if err != nil {
        self.defaultErrorReturn(err, c)
    }
    user, err := self.virtualRoomsManager.GetUser(session.Values["user_name"].(string))
    if err != nil {
        return c.Redirect(http.StatusSeeOther, "/app/user")
    }
    room, err := self.virtualRoomsManager.GetRoom(c.Param("roomName"))
    if err != nil {
        url := fmt.Sprintf("/app/user?error=%s", err.Error())
        return c.Redirect(303, url)
    }
        if room.Admin.Id != user.Id && room.Guest.Id != user.Id {
        url := fmt.Sprintf("/app/user?error=Not allowed in the room")
        return c.Redirect(303, url)
    }
    pageData["room"] = room
    pageData["user"] = user
    currUserId := session.Values["user_id"]
    if currUserId == room.Admin.Id {
        pageData["user_type"] = "Admin"
    } else {
        pageData["user_type"] = "Guest"
    }
    return c.Render(http.StatusOK,"room.tmpl", pageData)
}
func (self *VirtualRoomController) GuestRoomIndex(c echo.Context) error {
    pageData := self.getPageData(c)
    pageData["room_name"] = c.Param("roomName")
    pageData["user_type"] = "Guest"
    return c.Render(200, "guest.tmpl", pageData)
}
func (self *VirtualRoomController) GuestRoomJoin(c echo.Context) error {
    userName := c.FormValue("username")
    userExists, err := self.virtualRoomsManager.UserExists(userName)
    if err != nil {
        return self.defaultErrorReturn(err, c)
    }
    if userExists {
        return c.Redirect(303, "/join/room/"+c.Param("roomName")+"?error=User already exists")
    }
    user := virtualrooms.User {
        Id: uuid.New().String(),
        Name: userName,
    }
    if err := self.virtualRoomsManager.AddUser(user); err != nil {
        return self.defaultErrorReturn(err, c)
    }
    if err := self.virtualRoomsManager.GuestJoinRoom(user, c.FormValue("password"), c.Param("roomName")); err != nil {
        fmt.Println(err)
        self.virtualRoomsManager.DeleteKey("user:"+user.Name)
        return c.Redirect(303, "/join/room/"+c.Param("roomName")+"/?error="+err.Error())
    }
    session, err := self.getSession(c)
    session.Values["user_id"] = user.Id
    session.Values["user_name"] = user.Name
    if err != nil {
        fmt.Println(err)
        return c.Redirect(303, "/join/room/"+c.Param("roomName")+"/?error=" + err.Error())
    }
    //salvar sessao
    if err := session.Save(c.Request(), c.Response()); err != nil {
        fmt.Println(err)
        return c.Redirect(303, "/join/room/"+c.Param("roomName")+"/?error=" + err.Error())
    }
    return c.Redirect(303, "/app/room/"+c.Param("roomName"))
}
