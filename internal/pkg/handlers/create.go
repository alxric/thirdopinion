package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

// Create will render the create page
func Create(c echo.Context) error {
	return c.Render(http.StatusOK, "create", nil)
}
