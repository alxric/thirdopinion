package shared

import "myitcv.io/react"

//go:generate reactGen

// FooterDef is the definition for the Footer component
type FooterDef struct {
	react.ComponentDef
}

// FooterState is the state type for the Footer component
type FooterState struct {
}

// Footer creates instances of the Footer component
func Footer() *FooterElem {
	return buildFooterElem()
}

// Render renders the layout component
func (l FooterDef) Render() react.Element {
	return react.Div(
		&react.DivProps{
			ID: "footer",
		},
		react.S("FOOTER"),
	)
}
