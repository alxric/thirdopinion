package handlers

import (
	"errors"
	"net/http"
	"thirdopinion/internal/pkg/config"
	"thirdopinion/internal/pkg/psql"
	"time"

	"github.com/labstack/echo"
)

// ValidateSession against the DB
func ValidateSession(c echo.Context) error {
	resp := sessionResponse(c)
	if resp.Error != "" {
		c.Logger().Error(resp.Error)
	}
	return c.JSON(http.StatusOK, resp)
}

func sessionResponse(c echo.Context) *config.WSResponse {
	resp := &config.WSResponse{
		User: &config.User{},
	}
	headers := c.Request().Header
	if len(headers["Origin"]) != 1 || headers["Origin"][0] != "http://localhost:8080" {
		c.Logger().Error("Invalid origin header!")
		resp.Error = "Unauthorized request"
		return resp
	}
	var err error
	resp.User, err = validateToken(c)
	if err != nil {
		c.Logger().Error("No X-CSRF-TOKEN set!")
		resp.Error = "NO X-CSRF-TOKEN set!"
		return resp
	}
	err = psql.ValidateSession(resp.User)
	if err != nil {
		c.Logger().Error(err)
		resp.Error = "Invalid or expired token"
		go deleteCookie(c)
	}

	return resp
}

func validateToken(c echo.Context) (*config.User, error) {
	u := &config.User{}
	headers := c.Request().Header
	var token []string
	var ok bool
	if token, ok = headers["X-Csrf-Token"]; !ok {
		c.Logger().Error("No X-CSRF-TOKEN set!")
		return nil, errors.New("NO X-CSRF-TOKEN set")
	}
	u.SessionKey = token[0]
	return u, nil
}

func writeCookie(c echo.Context, r *config.WSResponse) error {
	cookie := new(http.Cookie)
	cookie.Name = "to_session_key"
	cookie.Value = r.User.SessionKey
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	c.SetCookie(cookie)
	return c.String(http.StatusOK, "")
}

func readCookie(c echo.Context) string {
	cookie, err := c.Cookie("to_session_key")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func deleteCookie(c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "to_session_key"
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(-24 * time.Hour)
	c.SetCookie(cookie)
}
