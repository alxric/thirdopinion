package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

// View will render the create page
func View(c echo.Context) error {
	return c.Render(http.StatusOK, "view", nil)
}
