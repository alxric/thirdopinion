package backend

import (
	"encoding/json"
	"thirdopinion/internal/pkg/config"
	"thirdopinion/internal/pkg/psql"
)

// Login sends the login request to the database
func Login(r *config.Register) ([]byte, error) {
	resp := &config.WSResponse{}

	u, err := psql.Login(r)
	switch err {
	case nil:
		err = psql.UpdateSession(u)
		if err != nil {
			return nil, err
		}
		resp.User = u
		resp.Msg = "logged in"
	default:
		resp.Error = err.Error()
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return b, nil
}
