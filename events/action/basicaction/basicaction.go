package basicaction

import "context"

// BasicAction is a basic implementation of the Action interface
type BasicAction func(context.Context)

// Run executes the action
func (a BasicAction) Run(ctx context.Context) {
	a(ctx)
}
