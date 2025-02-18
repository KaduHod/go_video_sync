package controllers

import (
	virtualrooms "kaduhod/video-sync/app/virtual_rooms"

	"github.com/labstack/echo/v4"
)
type SSEController struct {
   virtualRoomsManager *virtualrooms.VirtualRoomsManager
}
func NewSSEController(manager *virtualrooms.VirtualRoomsManager) *SSEController {
   return &SSEController{
      virtualRoomsManager: manager,
   }
}

func (self *SSEController) StreamRoom(c echo.Context) error {
    _ = c.Param("roomName")
    return nil
}

