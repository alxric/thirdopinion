package main

import (
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

// Layout creates instances of the Layout component
func Layout() *LayoutElem {
	return buildLayoutElem()
}

// Render renders the layout component
func (l LayoutDef) Render() react.Element {
	h := shared.HeaderDef{}
	f := shared.FooterDef{}
	return react.Fragment(
		h.Render(),
		react.Div(
			&react.DivProps{
				ClassName: "container-fluid",
			},
			react.Div(nil, Register()),
		),
		f.Render(),
	)
}
