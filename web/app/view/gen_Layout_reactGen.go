// Code generated by reactGen. DO NOT EDIT.

package main

import "myitcv.io/react"

type LayoutElem struct {
	react.Element
}

func buildLayout(cd react.ComponentDef) react.Component {
	return LayoutDef{ComponentDef: cd}
}

func buildLayoutElem(props LayoutProps, children ...react.Element) *LayoutElem {
	return &LayoutElem{
		Element: react.CreateElement(buildLayout, props, children...),
	}
}

func (l LayoutDef) RendersElement() react.Element {
	return l.Render()
}

// SetState is an auto-generated proxy proxy to update the state for the
// Layout component.  SetState does not immediately mutate l.State()
// but creates a pending state transition.
func (l LayoutDef) SetState(state LayoutState) {
	l.ComponentDef.SetState(state)
}

// State is an auto-generated proxy to return the current state in use for the
// render of the Layout component
func (l LayoutDef) State() LayoutState {
	return l.ComponentDef.State().(LayoutState)
}

// IsState is an auto-generated definition so that LayoutState implements
// the myitcv.io/react.State interface.
func (l LayoutState) IsState() {}

var _ react.State = LayoutState{}

// GetInitialStateIntf is an auto-generated proxy to GetInitialState
func (l LayoutDef) GetInitialStateIntf() react.State {
	return LayoutState{}
}

func (l LayoutState) EqualsIntf(val react.State) bool {
	return l == val.(LayoutState)
}

// IsProps is an auto-generated definition so that LayoutProps implements
// the myitcv.io/react.Props interface.
func (l LayoutProps) IsProps() {}

// Props is an auto-generated proxy to the current props of Layout
func (l LayoutDef) Props() LayoutProps {
	uprops := l.ComponentDef.Props()
	return uprops.(LayoutProps)
}

func (l LayoutProps) EqualsIntf(val react.Props) bool {
	return l.Equals(val.(LayoutProps))
}

var _ react.Props = LayoutProps{}