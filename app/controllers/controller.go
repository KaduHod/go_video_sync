package controllers

import (
	"fmt"
	"net/url"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)
type Controller struct {}
func (self *Controller) getPageData(c echo.Context) map[string]interface{} {
    pageData := make(map[string]interface{})
    if c.Get("csrf") != nil {
        pageData["csrf"] = c.Get("csrf").(string)
    }
    session, err := session.Get("session", c)
    if err != nil {
        fmt.Println(err)
        return pageData
    }
    if c.QueryParam("error") != "" {
        pageData["user_error"] = c.QueryParam("error")
    }
    pageData["sess"] = session.Values
    return pageData
}
func (self *Controller) getSession(c echo.Context) (*sessions.Session, error) {
    session, err := session.Get("session", c)
    if err != nil {
        return session, err
    }
    return session, nil
}
func (self Controller) getParam(paramName string, c echo.Context) string {
    url := c.Param(paramName)
    if url == "" {
        return url
    }
    return self.manipulaURL(url)
}
func (self Controller) manipulaURL(input string) (string) {
	if len(input) > 1 && input[len(input)-1] == '/' {
		input = input[:len(input)-1]
	}
	parsedURL, err := url.Parse(input)
	if err != nil {
		return ""
	}
	return parsedURL.Path
}
