package basicaction

import "context"

type BasicAction func(context.Context)

func (a BasicAction) Run(ctx context.Context) {
	a(ctx)
}
