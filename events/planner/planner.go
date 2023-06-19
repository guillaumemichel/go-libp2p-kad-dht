package planner

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/events"
)

type TimedAction any

// ActionPlanner is an interface for scheduling actions at a specific time.
type ActionPlanner interface {
	// ScheduleAction schedules an action to run at a specific time
	ScheduleAction(context.Context, time.Time, events.Action) TimedAction
	// RemoveAction removes an action from the planner
	RemoveAction(context.Context, TimedAction)

	// PopOverdueActions returns all actions that are overdue and removes them
	// from the planner
	PopOverdueActions(context.Context) []events.Action
}

// MultiActionPlanner is an interface for scheduling multiple actions at
// specific times.
type MultiActionPlanner interface {
	ActionPlanner

	// ScheduleActions schedules multiple actions at specific times
	ScheduleActions(context.Context, []time.Time, []events.Action)
	// RemoveActions removes multiple actions from the planner
	RemoveActions(context.Context, []events.Action)
}

// ScheduleActions schedules multiple actions at specific times using a planner.
func ScheduleActions(ctx context.Context, p ActionPlanner,
	times []time.Time, actions []events.Action) {

	switch p := p.(type) {
	case MultiActionPlanner:
		p.ScheduleActions(ctx, times, actions)
	default:
		for i, d := range times {
			p.ScheduleAction(ctx, d, actions[i])
		}
	}
}

// RemoveActions removes multiple actions from the planner.
func RemoveActions(ctx context.Context, p ActionPlanner, actions []events.Action) {
	switch p := p.(type) {
	case MultiActionPlanner:
		p.RemoveActions(ctx, actions)
	default:
		for _, a := range actions {
			p.RemoveAction(ctx, a)
		}
	}
}

// AwareActionPlanner is an interface for scheduling actions at a specific time
// and knowing when the next action will be scheduled.
type AwareActionPlanner interface {
	ActionPlanner

	// NextActionTime returns the time of the next action that will be
	// scheduled. If there are no actions scheduled, it returns MaxTime.
	NextActionTime(context.Context) time.Time
}
