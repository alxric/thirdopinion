package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"honnef.co/go/js/dom"
	"myitcv.io/react"
)

//go:generate reactGen

var (
	//emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// RegisterDef is the definition fo the Register component
type RegisterDef struct {
	react.ComponentDef
}

// RegisterState is the state type for the Register component
type RegisterState struct {
	email         string
	password      string
	confirm       string
	emailError    string
	passwordError string
	confirmError  string
}

// Register creates instances of the Register component
func Register() *RegisterElem {
	return buildRegisterElem()
}

// Equals must be defined because struct val instances of RegisterState cannot
// be compared. It is generally bad practice to have mutable values in state in
// this way; myitcv.io/immutable seeks to help address this problem.
// See myitcv.io/react/examples/immtodoapp for an example
func (c RegisterState) Equals(v RegisterState) bool {
	if c.email != v.email {
		return false
	}

	if c.password != v.password {
		return false
	}

	if c.confirm != v.confirm {
		return false
	}

	return true
}

// Render renders the Register component
func (t RegisterDef) Render() react.Element {
	submitClass := "rightFloat"
	switch {
	case t.State().email == "", t.State().password == "", t.State().confirm == "":
		submitClass += " invisible"
	case t.State().emailError != "", t.State().passwordError != "", t.State().confirmError != "":
		submitClass += " invisible"
	}
	return react.Fragment(
		react.Div(
			&react.DivProps{
				ID: "registerContainer",
			},
			react.Div(
				&react.DivProps{
					ID: "registerHeader",
				},
				react.H3(nil,
					react.S("Create a new account"),
				),
			),
			react.Form(
				&react.FormProps{
					ClassName: "form-inline",
				},
				react.Div(
					&react.DivProps{
						ClassName: "registerRow",
					},
					react.Div(
						&react.DivProps{
							ClassName: "form-group registerDiv",
							ID:        "registerEmailDiv",
						},
						react.Label(
							&react.LabelProps{
								ClassName: "sr-only",
								For:       "registerEmailInput",
							},
							react.S("Email input"),
						),
						react.Input(
							&react.InputProps{
								Type:        "text",
								ID:          "registerEmailInput",
								ClassName:   "form-control registerInput",
								Placeholder: "Enter email address",
								Value:       t.State().email,
								OnChange:    inputChange{t},
							},
						),
					),
					react.Div(
						&react.DivProps{
							ID:        "registerEmailError",
							ClassName: "error registerError",
						},
						react.S(t.State().emailError),
					),
				),
				react.Div(
					&react.DivProps{
						ClassName: "registerRow",
					},
					react.Div(
						&react.DivProps{
							ClassName: "form-group registerDiv",
							ID:        "registerPasswordDiv",
						},
						react.Label(
							&react.LabelProps{
								ClassName: "sr-only",
								For:       "registerPasswordInput",
							},
							react.S("Password input"),
						),
						react.Input(
							&react.InputProps{
								Type:        "password",
								ID:          "registerPasswordInput",
								ClassName:   "form-control registerInput",
								Placeholder: "Enter password",
								Value:       t.State().password,
								OnChange:    inputChange{t},
							},
						),
					),
					react.Div(
						&react.DivProps{
							ID:        "registerPasswordError",
							ClassName: "error registerError",
						},
						react.S(t.State().passwordError),
					),
				),
				react.Div(
					&react.DivProps{
						ClassName: "registerRow",
					},
					react.Div(
						&react.DivProps{
							ClassName: "form-group registerDiv",
							ID:        "confirmPasswordDiv",
						},
						react.Label(
							&react.LabelProps{
								ClassName: "sr-only",
								For:       "confirmPasswordInput",
							},
							react.S("Password input"),
						),
						react.Input(
							&react.InputProps{
								Type:        "password",
								ID:          "confirmPasswordInput",
								ClassName:   "form-control registerInput",
								Placeholder: "Confirm password",
								Value:       t.State().confirm,
								OnChange:    inputChange{t},
							},
						),
					),
					react.Div(
						&react.DivProps{
							ID:        "registerConfirmError",
							ClassName: "error registerError",
						},
						react.S(t.State().confirmError),
					),
				),
				react.Div(
					&react.DivProps{
						ID:        "registerSubmitDiv",
						ClassName: submitClass,
					},
					react.Button(
						&react.ButtonProps{
							Type:      "submit",
							ClassName: "btn btn-default",
							OnClick:   register{t},
						},
						react.S("Register"),
					),
				),
				react.Div(
					&react.DivProps{
						ClassName: "clear",
					},
				),
			),
		),
	)
}

type inputChange struct{ t RegisterDef }

func (i inputChange) OnChange(se *react.SyntheticEvent) {
	target := se.Target().(*dom.HTMLInputElement)
	ns := i.t.State()
	switch target.ID() {
	case "registerEmailInput":
		ns.emailError = ""
		ns.email = target.Value
		if ok := validateEmail(ns.email); !ok {
			ns.emailError = "Invalid email format"
		}
	case "registerPasswordInput":
		ns.passwordError = ""
		ns.confirmError = ""
		ns.password = target.Value
		if ok, msg := validatePassword(ns.password); !ok {
			ns.passwordError = msg
		}
	case "confirmPasswordInput":
		ns.confirmError = ""
		ns.confirm = target.Value
		if ns.confirm != ns.password {
			ns.confirmError = "Passwords don't match!"
		}
	}
	i.t.SetState(ns)
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

type register struct{ t RegisterDef }

func (r register) OnClick(se *react.SyntheticMouseEvent) {
	se.PreventDefault()
	ns := r.t.State()
	fmt.Println(ns)
}
