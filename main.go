package main

import (
	"fmt"
	"html/template"
	"io"
	"kaduhod/video-sync/app/controllers"
	app_middleware "kaduhod/video-sync/app/middlewares"
	virtualrooms "kaduhod/video-sync/app/virtual_rooms"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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
    e.Renderer = t
    e.Static("/public", "public")
    manager := virtualrooms.NewVirtualRoomsManager()
    manager.Init()
    vrmController := controllers.NewVirtualRoomController(manager)
    sseControler := controllers.NewSSEController(manager)
    e.GET("/", vrmController.Index)
    e.POST("/user", vrmController.NewUser)
    e.GET("/join/room/:roomName", vrmController.GuestRoomIndex)
    e.POST("/join/room/:roomName", vrmController.GuestRoomJoin)
    sessionGroup := e.Group("/app")
    sessionGroup.Use(app_middleware.OnlySessionInitialized)
    sessionGroup.GET("/user", vrmController.IndexUsers)
    sessionGroup.POST("/room/create", vrmController.NewRoom)
    sessionGroup.GET("/room/:roomName", vrmController.RoomIndex)
    sessionGroup.GET("/sse/:roomName", sseControler.StreamRoom)
    e.Logger.Fatal(e.Start(":3003"))
}
