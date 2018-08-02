package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

// Login will render the login page
func Login(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}
