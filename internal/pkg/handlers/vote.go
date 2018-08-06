package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"thirdopinion/internal/pkg/backend"
	"thirdopinion/internal/pkg/config"

	"github.com/labstack/echo"
)

// Vote for an argument
func Vote(c echo.Context) error {
	wsr := &config.WSRequest{}
	resp := &config.WSResponse{
		User: &config.User{},
	}
	if err := c.Bind(wsr); err != nil {
		resp.Error = err.Error()
		return c.JSON(http.StatusOK, resp)
	}

	resp = sessionResponse(c)
	wsr.Vote.User = resp.User.ID

	b, err := backend.Vote(wsr.Vote)
	if err != nil {
		resp.Error = err.Error()
		if strings.Contains(resp.Error, "violates unique constraint") {
			resp.Error = "You have already voted for this argument"
		}
		return c.JSON(http.StatusOK, resp)
	}
	if err := json.Unmarshal(b, resp); err != nil {
		resp.Error = err.Error()
	}
	return c.JSON(http.StatusOK, resp)
}
