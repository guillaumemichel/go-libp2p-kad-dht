package events

import (
	"context"
	"fmt"
)

type Action interface{}

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
