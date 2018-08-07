package handlers

import (
	"net/http"
	"thirdopinion/internal/pkg/backend"
	"thirdopinion/internal/pkg/config"

	"github.com/labstack/echo"
)

// Create will render the create page
func Create(c echo.Context) error {
	sessionKey := readCookie(c)
	return c.Render(http.StatusOK, "create", sessionKey)
}

// CreateArgument will add an argument to the database
func CreateArgument(c echo.Context) error {
	req := &config.WSRequest{}
	resp := &config.WSResponse{}
	if err := c.Bind(req); err != nil {
		resp.Error = err.Error()
		return c.JSON(http.StatusOK, resp)
	}
	resp = sessionResponse(c)
	if resp.Error != "" {
		return c.JSON(http.StatusOK, resp)
	}
	req.Argument.UserID = int64(resp.User.ID)
	resp2, err := backend.NewArgument(req.Argument)
	if err != nil {
		resp.Error = err.Error()
		return c.JSON(http.StatusOK, resp)
	}
	return c.JSON(http.StatusOK, resp2)
}
