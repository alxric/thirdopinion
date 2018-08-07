package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"thirdopinion/internal/pkg/backend"
	"thirdopinion/internal/pkg/config"

	"github.com/labstack/echo"
)

// Register will render the register page
func Register(c echo.Context) error {
	sessionKey := readCookie(c)
	return c.Render(http.StatusOK, "register", sessionKey)
}

// RegisterUser will create a new account
func RegisterUser(c echo.Context) error {
	wsr := &config.WSRequest{}
	if err := c.Bind(wsr); err != nil {
		return c.String(http.StatusForbidden, "Invalid request")
	}
	b, err := backend.Register(wsr.Register)
	if err != nil {
		return c.String(http.StatusForbidden, "Invalid request")
	}
	resp := &config.WSResponse{}
	if err := json.Unmarshal(b, resp); err != nil {
		return c.String(http.StatusForbidden, "Invalid backend response")
	}
	if strings.Contains(resp.Error, "violates unique constraint") {
		resp.Error = "Email address is already registered"
	}

	return c.JSON(http.StatusOK, resp)
}
