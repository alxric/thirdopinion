package backend

import (
	"encoding/json"
	"fmt"
	"strings"
	"thirdopinion/internal/pkg/config"
	"thirdopinion/internal/pkg/psql"
	"unicode"

	"github.com/gorilla/websocket"
)

// Register sends the register request to the database
func Register(ws *websocket.Conn, r *config.Register) ([]byte, error) {
	resp := &config.WSResponse{}
	pwOK, pwMsg := validatePassword(r.Password)
	switch {
	case r.Password != r.ConfirmPassword:
		resp.Error = "Passwords do not match"
	case !validateEmail(r.Email):
		resp.Error = "Invalid email address"
	case !pwOK:
		resp.Error = pwMsg
	}
	if resp.Error == "" {
		res, err := psql.Register(r)
		switch err {
		case nil:
			resp.Msg = res
		default:
			resp.Error = err.Error()
		}
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func validateEmail(email string) bool {
	if strings.Count(email, "@") != 1 {
		return false
	}
	vals := strings.Split(email, "@")
	if len(vals) != 2 || vals[0] == "" || vals[1] == "" {
		return false
	}
	return true
}

func validatePassword(pw string) (bool, string) {
	var uppercase, lowercase, digit, special bool
	if len(pw) < 8 || len(pw) > 16 {
		return false, "Password must be >= 8 <= 16 characters long"
	}
	for _, c := range pw {
		switch {
		case unicode.IsNumber(c):
			digit = true
		case unicode.IsUpper(c):
			uppercase = true
		case unicode.IsLower(c):
			lowercase = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		default:
			return false, fmt.Sprintf("Invalid character: %v", c)
		}
	}
	switch {
	case !uppercase:
		return false, "Password must contain at least one uppercase character"
	case !lowercase:
		return false, "Password must contain at least one lowercase character"
	case !digit:
		return false, "Password must contain at least one digit"
	case !special:
		return false, "Password must contain at least one special character"
	}
	return true, ""
}
