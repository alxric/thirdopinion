package main

import (
	"honnef.co/go/js/dom"
	"myitcv.io/react"
)

//go:generate reactGen

var document = dom.GetWindow().Document()

func main() {
	domTarget := document.GetElementByID("app")
	react.Render(Layout(), domTarget)
}
