package shared

import "myitcv.io/react"

//go:generate reactGen

// HeaderDef is the definition for the Header component
type HeaderDef struct {
	react.ComponentDef
}

// HeaderState is the state type for the Header component
type HeaderState struct {
}

// Header creates instances of the Header component
func Header() *HeaderElem {
	return buildHeaderElem()
}

// Render renders the layout component
func (l HeaderDef) Render() react.Element {
	return react.Div(
		&react.DivProps{
			ID: "header",
		},
		react.A(
			&react.AProps{
				Href: "/",
			},
			react.Img(
				&react.ImgProps{
					Src: "/static/img/logo.png",
					ID:  "logoImg",
				},
			),
		),
		react.S("Get a third opinion and solve your argument once and for all!"),
		react.Div(
			&react.DivProps{
				ClassName: "header-right",
			},
			react.S("Some buttons here"),
		),
	)
}
