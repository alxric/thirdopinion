package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"thirdopinion/internal/pkg/config"
	"time"

	"honnef.co/go/js/dom"
	"myitcv.io/react"
)

//go:generate reactGen

// CreateDef is the definition fo the Create component
type CreateDef struct {
	react.ComponentDef
}

// CreateState is the state type for the Create component
type CreateState struct {
	opinions     []*config.Opinion
	title        string
	nextPerson   int
	currItem     string
	fel          string
	err          *config.CreationResult
	lastUpdate   time.Time
	errClass     string
	initialClass string
	titleClass   string
}

// Create creates instances of the Create component
func Create() *CreateElem {
	return buildCreateElem()
}

// Equals must be defined because struct val instances of CreateState cannot
// be compared. It is generally bad practice to have mutable values in state in
// this way; myitcv.io/immutable seeks to help address this problem.
// See myitcv.io/react/examples/immtodoapp for an example
func (c CreateState) Equals(v CreateState) bool {
	if c.currItem != v.currItem {
		return false
	}

	if c.title != v.title {
		return false
	}

	if c.fel != v.fel {
		return false
	}

	if c.err != v.err {
		return false
	}

	if len(v.opinions) != len(c.opinions) {
		return false
	}

	for i := range v.opinions {
		if v.opinions[i] != c.opinions[i] {
			return false
		}
	}

	return true
}

// Render renders the Create component
func (t CreateDef) Render() react.Element {
	ph := "Enter the opinion of the first person"
	if len(t.State().opinions) >= 1 {
		ph = "Enter the opinion of the next person"
	}
	ns := t.handleErrors()
	var entries []react.Element
	entries = append(entries, t.generateEntries()...)
	if len(entries) > 0 {
		argumentHeader := []react.Element{react.Div(
			&react.DivProps{
				ID: "argumentHeader",
			},
			react.H3(nil, react.S("Arguments")),
		),
		}
		entries = append(argumentHeader, entries...)
	}
	return react.Fragment(
		react.Div(
			&react.DivProps{
				ID: "newArgumentDiv",
			},
			react.Form(&react.FormProps{ClassName: "form-inline"},
				react.Div(
					&react.DivProps{
						ID: "argumentTitle",
					},
					react.Div(
						&react.DivProps{
							ID: "argumentTitleText",
						},
						react.H3(nil, react.S("Title")),
					),
					react.Div(
						&react.DivProps{ClassName: "form-group", ID: "argumentTitleDiv"},
						react.Label(&react.LabelProps{ClassName: "sr-only", For: "argumentTitleInput"}, react.S("Title input")),
						react.Input(&react.InputProps{
							Type:        "text",
							ClassName:   ns.titleClass,
							ID:          "argumentTitleInput",
							Placeholder: "Enter argument title",
							Value:       t.State().title,
							OnChange:    titleChange{t},
						}),
					),
					react.Div(nil,
						react.Button(&react.ButtonProps{
							Type:      "submit",
							ClassName: "btn btn-default invisible",
							OnClick:   add{t},
						}, react.S("Add")),
					),
					react.Div(
						&react.DivProps{
							ID: "argumentInputs",
						},
						entries...,
					),
					react.Div(
						&react.DivProps{ClassName: "form-group", ID: "inputDiv1"},
						react.Label(&react.LabelProps{ClassName: "sr-only", For: "initialInput"}, react.S("Initial input")),
						react.Input(&react.InputProps{
							Type:        "text",
							ClassName:   ns.initialClass,
							ID:          "initialInput",
							Placeholder: ph,
							Value:       t.State().currItem,
							OnChange:    inputChange{t},
						}),
					),
					react.Div(nil,
						react.Div(
							&react.DivProps{
								ClassName: "leftFloat",
							},
							react.Span(
								&react.SpanProps{
									ClassName: ns.errClass,
									ID:        "errorSpan",
								},
								react.S(ns.err.Error),
							),
						),
						react.Div(
							&react.DivProps{
								ClassName: "rightFloat",
							},
							react.Button(
								&react.ButtonProps{
									Type:      "submit",
									ClassName: "btn btn-default",
									OnClick:   submit{t},
								},
								react.S("Submit"),
							),
						),
						react.Div(
							&react.DivProps{
								ClassName: "clear",
							},
						),
					),
				),
			),
		),
	)
}

func (t CreateDef) handleErrors() CreateState {
	ns := t.State()
	ns.titleClass = "form-control"
	ns.initialClass = "form-control"
	ns.errClass = "error"
	if ns.err == nil {
		ns.err = &config.CreationResult{}
	}
	switch ns.err.Position {
	case "title":
		ns.titleClass += " errorInput"
	case "argument":
		ns.initialClass += " errorInput"
	default:
		ns.errClass += " invisible"
	}
	return ns
}

func (t CreateDef) generateEntries() (entries []react.Element) {
	for index, opinion := range t.State().opinions {
		peopleSelect := react.Select(
			&react.SelectProps{
				ClassName: "peopleSelect form-control",
				Value:     fmt.Sprintf("Person %d", opinion.Person),
				ID:        fmt.Sprintf("peopleSelect%d", index),
				OnChange:  selectChange{t},
			},
			react.Option(nil, react.S("Person 1")),
			react.Option(nil, react.S("Person 2")),
		)
		entry :=
			react.Div(nil,
				react.Div(
					&react.DivProps{
						ClassName: "argumentInput",
					},
					react.Div(
						&react.DivProps{
							ClassName: "leftFloat peopleDiv",
						},
						react.Label(&react.LabelProps{ClassName: "sr-only", For: fmt.Sprintf("peopleSelect%d", index)}, react.S("Initial input")),
						peopleSelect,
					),
					react.Div(
						&react.DivProps{
							ClassName: "leftFloat argumentDiv",
						},
						react.S(opinion.Text),
					),
					react.Div(
						&react.DivProps{
							ClassName: "rightFloat deleteArgumentDiv",
						},
						react.Span(
							&react.SpanProps{
								ClassName: "deleteSpan",
								ID:        fmt.Sprintf("delete%d", index),
								OnClick:   del{t},
							},
							react.S("DELETE"),
						),
					),
					react.Div(
						&react.DivProps{
							ClassName: "clear",
						},
					),
				),
				react.Div(
					&react.DivProps{
						ClassName: "clear",
					},
				),
			)
		entries = append(entries, entry)
	}
	return

}

type inputChange struct{ t CreateDef }
type titleChange struct{ t CreateDef }
type add struct{ t CreateDef }
type del struct{ t CreateDef }
type submit struct{ t CreateDef }
type selectChange struct{ t CreateDef }

func (i inputChange) OnChange(se *react.SyntheticEvent) {
	target := se.Target().(*dom.HTMLInputElement)
	ns := i.t.State()
	ns.currItem = target.Value
	i.t.SetState(ns)
}

func (t titleChange) OnChange(se *react.SyntheticEvent) {
	target := se.Target().(*dom.HTMLInputElement)
	ns := t.t.State()
	ns.title = target.Value
	t.t.SetState(ns)
}

func (a add) OnClick(se *react.SyntheticMouseEvent) {
	ns := a.t.State()
	if ns.currItem != "" {
		person := 1
		if ns.nextPerson%2 != 0 {
			person = 2
		}
		opinion := &config.Opinion{
			Person: person,
			Text:   ns.currItem,
		}
		ns.opinions = append(ns.opinions, opinion)
		ns.currItem = ""
		ns.nextPerson++
	}
	a.t.SetState(ns)

	se.PreventDefault()
}

func (s submit) OnClick(se *react.SyntheticMouseEvent) {

	ns := s.t.State()
	arg := &config.Argument{
		Title:    ns.title,
		Opinions: ns.opinions,
	}
	se.PreventDefault()

	ch := writeDB(arg)
	go func() {
		select {
		case msg := <-ch:
			cr := &config.CreationResult{}
			err := json.Unmarshal([]byte(msg.Data.String()), cr)
			if err != nil {
				ns.err = &config.CreationResult{Error: "Unknown error"}
			}
			ns.err = cr
			s.t.SetState(ns)
			s.t.Render()
		}
	}()
}

func (d del) OnClick(se *react.SyntheticMouseEvent) {
	ns := d.t.State()
	target := se.Target().(*dom.HTMLSpanElement)
	id, err := strconv.Atoi(strings.Split(target.ID(), "delete")[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	ns.opinions = append(ns.opinions[:id], ns.opinions[id+1:]...)
	d.t.SetState(ns)
	d.t.ForceUpdate()
}

func (s selectChange) OnChange(se *react.SyntheticEvent) {
	ns := s.t.State()
	target := se.Target().(*dom.HTMLSelectElement)
	id, err := strconv.Atoi(strings.Split(target.ID(), "peopleSelect")[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	personNum, err := strconv.Atoi(strings.Split(target.Value, "Person ")[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	ns.opinions[id].Person = personNum
	s.t.SetState(ns)
	s.t.ForceUpdate()
}
