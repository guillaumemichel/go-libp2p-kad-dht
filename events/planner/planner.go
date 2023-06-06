package planner

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/events"
)

type ActionPlanner interface {
	ScheduleAction(context.Context, time.Time, events.Action)
	RemoveAction(context.Context, events.Action)
	PopOverdueActions(context.Context) []events.Action
}

type MultActionPlanner interface {
	ActionPlanner
	ScheduleActions(context.Context, []time.Time, []events.Action)
	RemoveActions(context.Context, []events.Action)
}

func ScheduleActions(ctx context.Context, p ActionPlanner,
	times []time.Time, actions []events.Action) {

	switch p := p.(type) {
	case MultActionPlanner:
		p.ScheduleActions(ctx, times, actions)
	default:
		for i, d := range times {
			p.ScheduleAction(ctx, d, actions[i])
		}
	}
}

func RemoveActions(ctx context.Context, p ActionPlanner, actions []events.Action) {
	switch p := p.(type) {
	case MultActionPlanner:
		p.RemoveActions(ctx, actions)
	default:
		for _, a := range actions {
			p.RemoveAction(ctx, a)
		}
	}
}
