package main

import (
	"fmt"
	"thirdopinion/internal/pkg/config"
	"thirdopinion/web/app/shared"

	"honnef.co/go/js/dom"
	"myitcv.io/react"
)

//go:generate reactGen

// LoginDef is the definition fo the Login component
type LoginDef struct {
	react.ComponentDef
}

// LoginState is the state type for the Login component
type LoginState struct {
	email      string
	password   string
	loginError string
}

// Login creates instances of the Login component
func Login() *LoginElem {
	return buildLoginElem()
}

// Equals must be defined because struct val instances of LoginState cannot
// be compared. It is generally bad practice to have mutable values in state in
// this way; myitcv.io/immutable seeks to help address this problem.
// See myitcv.io/react/examples/immtodoapp for an example
func (c LoginState) Equals(v LoginState) bool {
	if c.email != v.email {
		return false
	}

	if c.password != v.password {
		return false
	}

	if c.loginError != v.loginError {
		return false
	}

	return true
}

// Render renders the Login component
func (t LoginDef) Render() react.Element {
	return react.Fragment(
		react.Div(
			&react.DivProps{
				ID: "loginContainer",
			},
			react.Div(
				&react.DivProps{
					ID: "loginHeader",
				},
				react.H3(nil,
					react.S("Login to your account"),
				),
			),
			react.Form(
				&react.FormProps{
					ClassName: "form-inline",
				},
				react.Div(
					&react.DivProps{
						ClassName: "loginRow",
					},
					react.Div(
						&react.DivProps{
							ClassName: "form-group loginDiv",
							ID:        "loginEmailDiv",
						},
						react.Label(
							&react.LabelProps{
								ClassName: "sr-only",
								For:       "loginEmailInput",
							},
							react.S("Email input"),
						),
						react.Input(
							&react.InputProps{
								Type:        "text",
								ID:          "loginEmailInput",
								ClassName:   "form-control loginInput",
								Placeholder: "Enter email address",
								Value:       t.State().email,
								OnChange:    inputChange{t},
							},
						),
					),
				),
				react.Div(
					&react.DivProps{
						ClassName: "loginRow",
					},
					react.Div(
						&react.DivProps{
							ClassName: "form-group loginDiv",
							ID:        "loginPasswordDiv",
						},
						react.Label(
							&react.LabelProps{
								ClassName: "sr-only",
								For:       "loginPasswordInput",
							},
							react.S("Password input"),
						),
						react.Input(
							&react.InputProps{
								Type:        "password",
								ID:          "loginPasswordInput",
								ClassName:   "form-control loginInput",
								Placeholder: "Enter password",
								Value:       t.State().password,
								OnChange:    inputChange{t},
							},
						),
					),
					react.Div(
						&react.DivProps{
							ID:        "loginError",
							ClassName: "error loginError",
						},
						react.S(t.State().loginError),
					),
				),
				react.Div(
					&react.DivProps{
						ID:        "loginSubmitDiv",
						ClassName: "rightFloat",
					},
					react.Button(
						&react.ButtonProps{
							Type:      "submit",
							ClassName: "btn btn-default",
							OnClick:   login{t},
						},
						react.S("Login"),
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

type inputChange struct{ t LoginDef }

func (i inputChange) OnChange(se *react.SyntheticEvent) {
	target := se.Target().(*dom.HTMLInputElement)
	ns := i.t.State()
	switch target.ID() {
	case "loginEmailInput":
		ns.email = target.Value
	case "loginPasswordInput":
		ns.loginError = ""
		ns.password = target.Value
	}
	i.t.SetState(ns)
}

type login struct{ t LoginDef }

func (l login) OnClick(se *react.SyntheticMouseEvent) {
	ns := l.t.State()
	wsr := &config.WSRequest{
		Method: "Login",
		Register: &config.Register{
			Email:    ns.email,
			Password: ns.password,
		},
	}
	se.PreventDefault()
	ch := shared.WriteDB(wsr)

	go func() {
		select {
		case wr := <-ch:
			if wr.Error != "" {
				ns.loginError = wr.Error
			}
			if wr.Msg == "Logged in" {
				fmt.Println("logged in success!!!!")
			}
			l.t.SetState(ns)
			l.t.ForceUpdate()
			l.t.Render()
		}
	}()
}
