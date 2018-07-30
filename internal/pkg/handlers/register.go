package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

// Register will render the register page
func Register(c echo.Context) error {
	return c.Render(http.StatusOK, "register", nil)
}
