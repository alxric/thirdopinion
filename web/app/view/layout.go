// layout.go

package main

import (
	"thirdopinion/internal/pkg/config"
	"thirdopinion/web/app/shared"

	"myitcv.io/react"
)

//go:generate reactGen

// LayoutDef is the definition for the Layout component
type LayoutDef struct {
	react.ComponentDef
}

// LayoutState is the state type for the Layout component
type LayoutState struct {
}

// LayoutProps contains the initial elements needed to generate a view
type LayoutProps struct {
	arguments []*config.Argument
}

// Layout creates instances of the Layout component
func Layout(p LayoutProps) *LayoutElem {
	return buildLayoutElem(p)
}

// Equals must be defined because struct val instances of LayoutState cannot
// be compared. It is generally bad practice to have mutable values in state in
// this way; myitcv.io/immutable seeks to help address this problem.
// See myitcv.io/react/examples/immtodoapp for an example
func (p LayoutProps) Equals(v LayoutProps) bool {
	return true
}

// Render renders the layout component
func (l LayoutDef) Render() react.Element {
	h := shared.HeaderDef{}
	f := shared.FooterDef{}
	vp := ViewProps{
		arguments: l.Props().arguments,
	}
	return react.Fragment(
		h.Render(),
		react.Div(
			&react.DivProps{
				ClassName: "container-fluid",
			},
			react.Div(nil, View(vp)),
		),
		f.Render(),
	)
}
