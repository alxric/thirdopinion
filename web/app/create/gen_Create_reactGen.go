// Code generated by reactGen. DO NOT EDIT.

package main

import "myitcv.io/react"

type CreateElem struct {
	react.Element
}

func buildCreate(cd react.ComponentDef) react.Component {
	return CreateDef{ComponentDef: cd}
}

func buildCreateElem(children ...react.Element) *CreateElem {
	return &CreateElem{
		Element: react.CreateElement(buildCreate, nil, children...),
	}
}

func (c CreateDef) RendersElement() react.Element {
	return c.Render()
}

// SetState is an auto-generated proxy proxy to update the state for the
// Create component.  SetState does not immediately mutate c.State()
// but creates a pending state transition.
func (c CreateDef) SetState(state CreateState) {
	c.ComponentDef.SetState(state)
}

// State is an auto-generated proxy to return the current state in use for the
// render of the Create component
func (c CreateDef) State() CreateState {
	return c.ComponentDef.State().(CreateState)
}

// IsState is an auto-generated definition so that CreateState implements
// the myitcv.io/react.State interface.
func (c CreateState) IsState() {}

var _ react.State = CreateState{}

// GetInitialStateIntf is an auto-generated proxy to GetInitialState
func (c CreateDef) GetInitialStateIntf() react.State {
	return CreateState{}
}

func (c CreateState) EqualsIntf(val react.State) bool {
	return c.Equals(val.(CreateState))
}
