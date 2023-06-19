package testaction

import "context"

type IntAction int

func (a IntAction) Run(context.Context) {}

type FuncAction struct {
	Ran bool
	Int int
}

func NewFuncAction(i int) *FuncAction {
	return &FuncAction{Int: i}
}

func (a *FuncAction) Run(context.Context) {
	a.Ran = true
}
