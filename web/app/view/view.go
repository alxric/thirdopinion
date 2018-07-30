package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
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
	for index, arg := range t.Props().arguments {
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
			t.generateVoteArea(arg, index),
		)
		arguments = append(arguments, argument)
	}
	return
}

func (t ViewDef) generateVoteArea(arg *config.Argument, index int) (voteArea react.Element) {
	totalVotes := float64(arg.Votes.Person1) + float64(arg.Votes.Person2)
	voteArea = react.Div(
		&react.DivProps{
			ClassName: "voteWrapper",
		},
		react.Div(
			&react.DivProps{
				ClassName: "voteHeader",
			},
			react.S("Votes so far"),
		),
		react.Div(
			&react.DivProps{
				ClassName: "voteDisplay",
			},
			react.Span(
				&react.SpanProps{
					ClassName: "person1 voteSpan",
				},
				react.S(fmt.Sprintf("%d (%d%%)", arg.Votes.Person1,
					int64(math.Round(100*(float64(arg.Votes.Person1)/totalVotes)))),
				),
			),
			react.Span(
				&react.SpanProps{
					ClassName: "person2 voteSpaG",
				},
				react.S(fmt.Sprintf("%d (%d%%)", arg.Votes.Person2,
					int64(math.Round(100*(float64(arg.Votes.Person2)/totalVotes)))),
				),
			),
		),
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
					ID:        fmt.Sprintf("voteButton_%d_%d_1", index, arg.ID),
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
					ID:        fmt.Sprintf("voteButton_%d_%d_2", index, arg.ID),
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

	bv, err := parseButtonID(target.ID())
	if err != nil {
		fmt.Println("Invalid button ID")
		return
	}

	ch, err := voteDB(bv.ArgID, bv.Person)
	if err != nil {
		fmt.Println("Error occured. Fix this alex!!!!")
		fmt.Println(err)
	}
	go func() {
		select {
		case msg := <-ch:
			wsr := &config.WSResponse{}
			err := json.Unmarshal([]byte(msg.Data.String()), wsr)
			if err != nil {
				fmt.Println("Could not unmarshal wsr for vote. This should genereate a proper error Alex!!")
			}
			switch wsr.Error {
			case "":
				switch wsr.Msg {
				case "Voted":
					switch bv.Person {
					case 1:
						v.t.Props().arguments[bv.Index].Votes.Person1++
					case 2:
						v.t.Props().arguments[bv.Index].Votes.Person2++
					}
					v.t.ForceUpdate()
					v.t.Render()
				}
			default:
				fmt.Println("WSR contains error! This should generate a proper error Alex!!!")
			}
		}
	}()
}

type buttonVals struct {
	Index  int
	ArgID  int
	Person int
}

func parseButtonID(id string) (*buttonVals, error) {
	idVals := strings.Split(id, "_")
	if len(idVals) != 4 {
		return nil, fmt.Errorf("Invalid button ID")
	}
	iIndex, err := strconv.Atoi(idVals[1])
	if err != nil {
		return nil, err
	}
	iArgID, err := strconv.Atoi(idVals[2])
	if err != nil {
		return nil, err
	}
	iPerson, err := strconv.Atoi(idVals[3])
	if err != nil {
		return nil, err
	}
	bv := &buttonVals{
		Index:  iIndex,
		ArgID:  iArgID,
		Person: iPerson,
	}
	return bv, nil

}
