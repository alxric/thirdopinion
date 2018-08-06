package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"thirdopinion/internal/pkg/backend"
	"thirdopinion/internal/pkg/config"

	"github.com/labstack/echo"
)

type viewVars struct {
	SessionKey string
	Votes      string
}

// View will render the view page
func View(c echo.Context) error {
	sessionKey := readCookie(c)
	votes, err := backend.UserVotes(sessionKey)
	if err != nil {
		c.Logger().Error(err)
	}
	jsonVotes, err := json.Marshal(votes)
	if err != nil {
		c.Logger().Error(err)
	}
	vv := viewVars{
		SessionKey: sessionKey,
		Votes:      string(jsonVotes),
	}
	return c.Render(http.StatusOK, "view", vv)
}

// ListArguments is the api endpoint for listing arguments
func ListArguments(c echo.Context) error {
	wsr := &config.WSRequest{
		Argument: &config.Argument{},
	}
	argID := c.QueryParam("id")
	if iArgID, err := strconv.Atoi(argID); err == nil {
		wsr.Argument.ID = int64(iArgID)
	}

	b, err := backend.ListArguments(wsr.Argument)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	resp := &config.WSResponse{}
	args := []*config.Argument{}
	err = json.Unmarshal(b, &args)
	if err != nil {
		return err
	}
	resp.Arguments = args
	if len(resp.Arguments) == 0 {
		resp.Error = "Post not found"
	}
	return c.JSON(http.StatusOK, resp)
}
