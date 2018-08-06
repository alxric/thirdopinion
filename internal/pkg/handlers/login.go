package handlers

import (
	"encoding/json"
	"net/http"
	"thirdopinion/internal/pkg/backend"
	"thirdopinion/internal/pkg/config"

	"github.com/labstack/echo"
)

// Login will render the login page
func Login(c echo.Context) error {
	sessionKey := readCookie(c)
	return c.Render(http.StatusOK, "login", sessionKey)
}

// NewLogin will register a new login
func NewLogin(c echo.Context) error {
	wsr := &config.WSRequest{}
	resp := &config.WSResponse{}
	if err := c.Bind(wsr); err != nil {
		resp.Error = err.Error()
		return c.JSON(http.StatusOK, resp)
	}
	b, err := backend.Login(wsr.Register)
	if err != nil {
		resp.Error = err.Error()
		return c.JSON(http.StatusOK, resp)
	}
	if err := json.Unmarshal(b, resp); err != nil {
		return err
	}
	if resp.Error != "" {
		return c.JSON(http.StatusOK, resp)
	}
	err = writeCookie(c, resp)
	if err != nil {
		resp.Error = err.Error()
		return c.JSON(http.StatusNoContent, resp)
	}
	return c.JSON(http.StatusOK, resp)
}
