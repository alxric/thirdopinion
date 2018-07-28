package main

import (
	"fmt"
	"thirdopinion/internal/pkg/config"

	"honnef.co/go/js/dom"
	"myitcv.io/react"
)

//go:generate reactGen

// ViewDef is the definition fo the Create component
type ViewDef struct {
	react.ComponentDef
}

// ViewState is the state type for the View component
type ViewState struct {
	err *config.CreationResult
}

// ViewProps contains the initial elements needed to generate a view
type ViewProps struct {
	arguments []*config.Argument
	err       *config.CreationResult
}

// View creates instances of the View component
func View(p ViewProps) *ViewElem {
	return buildViewElem(p)
}

// Equals must be defined because struct val instances of ViewState cannot
// be compared. It is generally bad practice to have mutable values in state in
// this way; myitcv.io/immutable seeks to help address this problem.
// See myitcv.io/react/examples/immtodoapp for an example
func (p ViewProps) Equals(v ViewProps) bool {

	if p.err != v.err {
		return false
	}

	return true
}

// Render renders the View component
func (t ViewDef) Render() react.Element {
	var arguments []react.Element
	arguments = append(arguments, t.generateArguments()...)
	return react.Fragment(
		react.Div(
			&react.DivProps{
				ID: "arguments",
			},
			arguments...,
		),
	)
}

func (t ViewDef) generateArguments() (arguments []react.Element) {
	for _, arg := range t.Props().arguments {
		var opinions []react.Element
		for _, opinion := range arg.Opinions {
			opDiv := react.Div(
				&react.DivProps{
					ClassName: "opinionWrapper",
				},
				react.Div(
					&react.DivProps{
						ClassName: fmt.Sprintf("person%d", opinion.Person),
					},
					react.S(fmt.Sprintf("<Person %d> %s", opinion.Person, opinion.Text)),
				),
			)
			opinions = append(opinions, opDiv)
		}
		argument := react.Div(
			&react.DivProps{
				ClassName: "argumentWrapper",
			},
			react.Div(
				&react.DivProps{
					ClassName: "argumentHeader",
				},
				react.A(
					&react.AProps{
						ClassName: "argumentLink",
						Href:      fmt.Sprintf("/view/%d", arg.ID),
					},
					react.H3(
						nil,
						react.S(arg.Title),
					),
				),
			),
			react.Div(
				&react.DivProps{
					ClassName: "opinionContainer",
				},
				opinions...,
			),
			t.generateVoteArea(arg),
		)
		arguments = append(arguments, argument)
	}
	return
}

func (t ViewDef) generateVoteArea(arg *config.Argument) (voteArea react.Element) {
	voteArea = react.Div(
		&react.DivProps{
			ClassName: "voteWrapper",
		},
		react.Div(
			&react.DivProps{
				ClassName: "voteHeader",
			},
			react.S("Who do you agree with?"),
		),
		react.Div(
			&react.DivProps{
				ClassName: "voteButtonDiv leftFloat",
			},
			react.Button(
				&react.ButtonProps{
					Type:      "submit",
					ID:        fmt.Sprintf("voteButton_%d_1", arg.ID),
					ClassName: "btn btn-default person1",
					OnClick:   vote{t},
				},
				react.S("Person 1"),
			),
		),
		react.Div(
			&react.DivProps{
				ClassName: "voteButtonDiv leftFloat",
			},
			react.Button(
				&react.ButtonProps{
					Type:      "submit",
					ID:        fmt.Sprintf("voteButton_%d_2", arg.ID),
					ClassName: "btn btn-default person2",
					OnClick:   vote{t},
				},
				react.S("Person 2"),
			),
		),
		react.Div(
			&react.DivProps{
				ClassName: "clear",
			},
		),
	)
	return
}

type vote struct{ t ViewDef }

func (v vote) OnClick(se *react.SyntheticMouseEvent) {
	target := se.Target().(*dom.HTMLButtonElement)
	se.PreventDefault()

	ch, err := voteDB(target.ID())
	if err != nil {
		fmt.Println("Error occured. Fix this alex!!!!")
		fmt.Println(err)
	}
	go func() {
		select {
		case msg := <-ch:
			fmt.Println(msg.Data.String())
		}
	}()
}
