package events

import (
	"context"
	"fmt"
)

// Action usually represents a function executed by a peer's scheduler.
// TODO: define the exact type of Action
type Action any

func Run(ctx context.Context, a Action) {
	switch a := a.(type) {
	case func():
		a()
	case func(context.Context):
		a(ctx)
	default:
		panic(fmt.Sprintf("unknown action type: %T", a))
	}
}
