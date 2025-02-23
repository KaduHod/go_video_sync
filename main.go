package main

import (
	"fmt"
	"html/template"
	"io"
	"kaduhod/video-sync/app/controllers"
	app_middleware "kaduhod/video-sync/app/middlewares"
	virtualrooms "kaduhod/video-sync/app/virtual_rooms"
	"net/http"
	"time"

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
	//c.Logger().Error(err)
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.File(errorPage); err != nil {
		//c.Logger().Error(err)
	}
    if code == 404 {
        c.Redirect(303, "/")
    }
}
func main() {
    e := echo.New()
    //e.HTTPErrorHandler = customHTTPErrorHandler
    t := &Template{
        templates: template.Must(template.ParseGlob("./views/*.tmpl")),
    }
	currentTime := func() string {
		return time.Now().Format("02/01/2006 15:04:05.000") // Dia/mês/ano Hora:minuto:segundo.milissegundos
	}
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
        Format: fmt.Sprintf("[%s] path: ${uri} ${method} | status: ${status} | lat: ${latency}\n", currentTime()),
	}))
    var store = sessions.NewCookieStore([]byte("Chave aleatória"))
    e.Use(session.MiddlewareWithConfig(session.Config{
        Store: store,
    }))
    e.Renderer = t
    e.Static("/public", "public")
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"http://localhost:3003"},
    }))
    manager := virtualrooms.NewVirtualRoomsManager()
    manager.Init()
    sseManager := virtualrooms.NewSSEManager(manager)
    vrmController := controllers.NewVirtualRoomController(manager)
    sseControler := controllers.NewSSEController(manager, sseManager)
    notSessionGroup := e.Group("")
    notSessionGroup.Use(app_middleware.OnlySessionNotInitialized)
    notSessionGroup.GET("/", vrmController.Index)
    notSessionGroup.POST("/user", vrmController.NewUser)
    notSessionGroup.GET("/join/room/:roomName", vrmController.GuestRoomIndex)
    notSessionGroup.POST("/join/room/:roomName", vrmController.GuestRoomJoin)
    sessionGroup := e.Group("/app")
    sessionGroup.Use(app_middleware.OnlySessionInitialized)
    sessionGroup.GET("/user", vrmController.IndexUsers)
    sessionGroup.POST("/room/create", vrmController.NewRoom)
    sessionGroup.GET("/room/:roomName", vrmController.RoomIndex)
    sessionGroup.POST("/room/:roomName/send", sseControler.Post)
    sessionGroup.DELETE("/room/close/:roomName", sseControler.CloseRoom)
    sessionGroup.GET("/join/room/:roomName", vrmController.LoguedGuestRoomJoinForm)
    sessionGroup.POST("/join/room/:roomName", vrmController.LoguedGuestRoomJoin)
    e.GET("/sse/room/:roomName", sseControler.StreamRoom)
    for _, route := range e.Routes() {
        fmt.Printf("Method: %s, Path: %s\n", route.Method, route.Path)
    }
    e.Logger.Fatal(e.Start(":3003"))
}
