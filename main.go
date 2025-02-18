package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	app_middleware "kaduhod/video-sync/app/middlewares"
	"kaduhod/video-sync/app/utils"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type VideoUpdate struct {
    Action string `json:"action"`
    SenderId string `json:"sender_id"`
    RoomId int `json:"room_id"`
}
type User struct {
    Id int `json:"id"`
    Name string `json:"name"`
}
type Room struct {
    Id int `json:"id"`
    Name string `json:"name"`
    Admin User `json:"admin"`
    Guest User `json:"guest"`
    Password string `json:"password"`
    GuestLink template.URL `json:"guest_link"`
}
/*
* Acoes de video
* Play
* Pause
* Forward
* Backward
*/
type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
func customHTTPErrorHandler(err error, c echo.Context) {
 	if c.Response().Committed {
 		return
 	}

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
    if code == 404 {
        c.Redirect(303, "/")
    }
}
func GetPageData(c echo.Context) map[string]interface{} {
    pageData := make(map[string]interface{})
    pageData["csrf"] = c.Get("csrf").(string)
    session, err := session.Get("session", c)
    if err != nil {
        return pageData
    }
    if c.QueryParam("error") != "" {
        pageData["user_error"] = c.QueryParam("error")
    }
    pageData["sess"] = session.Values
    return pageData
}
func main() {
    e := echo.New()
    e.HTTPErrorHandler = customHTTPErrorHandler
    t := &Template{
        templates: template.Must(template.ParseGlob("./views/*.tmpl")),
    }
    e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
        TokenLookup: "form:csrf",
    }))
    var store = sessions.NewCookieStore([]byte("Chave aleatória"))
    e.Use(session.MiddlewareWithConfig(session.Config{
        Store: store,
    }))
    var rooms []Room
    var users []User
    userAdmin := User {
        Id: 1,
        Name: "Admin",
    }
    hash, _ := utils.HashPassword("123456")
    defaultRoom := Room {
        Id: 1,
        Name: "Teste",
        Admin: userAdmin,
        Password: hash,
        GuestLink: template.URL("http://localhost:3003/join/room/Teste"),
    }
    usersLimit := 10
    roomsLimit := usersLimit/2
    rooms = append(rooms, defaultRoom)
    users = append(users, userAdmin)
    e.Renderer = t
    e.Static("/public", "public")
    sessionGroup := e.Group("/app")
    sessionGroup.Use(app_middleware.OnlySessionInitialized)
    e.GET("/", func(c echo.Context) error {
        pageData := GetPageData(c)
        if c.QueryParam("error") != "" {
            pageData["user_error"] = c.QueryParam("error")
        }
        return c.Render(200, "main.tmpl", pageData)
    })
    sessionGroup.GET("/user", func(c echo.Context) error {
        pageData := GetPageData(c)
        error := c.QueryParam("error")
        if error != "" {
            pageData["user_error"] = error
            return c.Render(400, "user.tmpl", pageData)
        }
        return c.Render(200, "user.tmpl", pageData)
    })
    e.POST("/user", func(c echo.Context) error {
        if len(users) > usersLimit {
            return c.Redirect(303, "/?error=Too many users")
        }
        name := c.FormValue("name")
        for _, user := range users {
            if user.Name == name {
                return c.Redirect(303, "/?error=User already exists")
            }
        }
        id := len(users) + 1
        user := User{
            Name: name,
            Id: id,
        }
        sess, err := session.Get("session", c)
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
        users = append(users, user)
        fmt.Println(users)
        return c.Redirect(303, "/app/user")
    })
    e.GET("/join/room/:roomName", func(c echo.Context) error {
        pageData := GetPageData(c)
        pageData["room_name"] = c.Param("roomName")
        return c.Render(200, "guest.tmpl", pageData)
    })
    e.POST("/join/room/:roomName", func(c echo.Context) error {
        // verificar se nome de usuario já existe
        name := c.FormValue("username")
        for _, user := range users {
            if user.Name == name {
                return c.Redirect(303, "/join/room/"+c.Param("roomName")+"?error=User already exists")
            }
        }
        // verificar se sala existe
        roomName := c.Param("roomName")
        exists := false
        for _, room := range rooms {
            if room.Name == roomName {
                exists = true
                break;
            }
        }
        if !exists {
            return c.Redirect(303, "/?error=Room not found")
        }
        // verificar senha
        room, err := getRoomByName(roomName, rooms)
        if err != nil {
            return c.Redirect(303, "/?error=" + err.Error())
        }
        if !utils.CheckPasswordHash(c.FormValue("password"), room.Password) {
            return c.Redirect(303, "/join/room/"+c.Param("roomName")+"?error=Wrong password")
        }
        // criar ususario
        user := User {
            Id: len(users) + 1,
            Name: name,
        }
        if room.Guest.Id != 0 {
            return c.Redirect(303, "/?error=Room is full")
        }
        session, err := session.Get("session", c)
        if err != nil {
            return c.Redirect(303, "/join/room/"+c.Param("roomName")+"/?error=" + err.Error())
        }
        // criar sessao
        session.Values["user_id"] = user.Id
        session.Values["user_name"] = user.Name
        users = append(users, user)
        room.Guest = user
        // Adicionar usuario como guest
        for key, room := range rooms {
            if room.Name == c.Param("roomName") {
                rooms[key].Guest = user
                break;
            }
        }
        return c.Redirect(303, "/app/room/"+c.Param("roomName"))
    })
    sessionGroup.POST("/room/create", func(c echo.Context) error {
        pageData := GetPageData(c)
        if len(rooms) > roomsLimit {
            pageData["user_error"] = "Too many rooms"
            return c.Render(400, "main.tmpl", pageData)
        }
        name := c.FormValue("name")
        for _, room := range rooms {
            if room.Name == name {
                return c.Redirect(303, "/app/user?error=Room name unavailable")
            }
        }
        user, err := getUser(users, c)
        if err != nil {
            return c.Redirect(303, "/app/user?error=" + err.Error())
        }
        password := c.FormValue("password")
        senhaHash, err := utils.HashPassword(password)
        if err != nil {
            return c.Redirect(303, "/app/user?error=" + err.Error())
        }
        room := Room{
            Id: len(rooms) + 1,
            Name: name,
            Admin: user,
            Password: senhaHash,
            GuestLink: template.URL(c.Request().Host+"/join/room/" + name),
        }
        rooms = append(rooms, room)
        roomUrl := fmt.Sprintf("/app/room/%s", name)
        return c.Redirect(http.StatusSeeOther, roomUrl)
    })
    sessionGroup.GET("/room/:roomname", func(c echo.Context) error {
        pageData := GetPageData(c)
        user, err := getUser(users, c)
        if err != nil {
            return c.Redirect(http.StatusSeeOther, "/app/user")
        }
        room, err := getRoomByName(c.Param("roomname"), rooms)
        if err != nil {
            url := fmt.Sprintf("/app/user?error=%s", err.Error())
            return c.Redirect(303, url)
        }
        if room.Admin.Id != user.Id && room.Guest.Id != user.Id {
            url := fmt.Sprintf("/app/user?error=Not allowed in the room")
            return c.Redirect(303, url)
        }
        pageData["room"] = room
        return c.Render(http.StatusOK,"room.tmpl", pageData)
    })
    e.Logger.Fatal(e.Start(":3003"))

}

func getUser(usersList []User, c echo.Context) (User, error) {
    var user User
    session, err := session.Get("session", c)
    if err != nil {
        return user, err
    }
    sessionUserId := session.Values["user_id"].(int)
    for _, userCurr := range usersList {
        if  sessionUserId == userCurr.Id {
            user = userCurr
            break
        }
    }
    if user.Id == 0 {
        return user, errors.New("User not found")
    }
    return user, nil
}
func getRoomByName(roomName string, roomsList []Room) (Room, error) {
    var room Room
    for _, roomCurr := range roomsList {
        if roomCurr.Name == roomName {
            room = roomCurr
            break
        }
    }
    if room.Id == 0 {
        return room, errors.New("Room not found")
    }
    return room, nil
}
