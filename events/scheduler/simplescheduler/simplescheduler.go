package simplescheduler

import (
	"context"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p-kad-dht/events"
	"github.com/libp2p/go-libp2p-kad-dht/events/planner"
	sp "github.com/libp2p/go-libp2p-kad-dht/events/planner/simpleplanner"
	"github.com/libp2p/go-libp2p-kad-dht/events/queue"
	"github.com/libp2p/go-libp2p-kad-dht/events/queue/chanqueue"
)

// SimpleScheduler is a simple implementation of the Scheduler interface. It
// uses a simple planner and a channel-based queue.
type SimpleScheduler struct {
	clk clock.Clock

	queue   queue.EventQueue
	planner planner.AwareActionPlanner
}

// NewSimpleScheduler creates a new SimpleScheduler.
func NewSimpleScheduler(ctx context.Context, clk clock.Clock) *SimpleScheduler {
	return &SimpleScheduler{
		clk: clk,

		queue:   chanqueue.NewChanQueue(ctx, 100),
		planner: sp.NewSimplePlanner(clk),
	}
}

// Now returns the scheduler's current time.
func (s *SimpleScheduler) Now() time.Time {
	return s.clk.Now()
}

// EnqueueAction enqueues an action to be run as soon as possible.
func (s *SimpleScheduler) EnqueueAction(ctx context.Context, a events.Action) {
	s.queue.Enqueue(ctx, a)
}

// ScheduleAction schedules an action to run at a specific time.
func (s *SimpleScheduler) ScheduleAction(ctx context.Context, t time.Time, a events.Action) {
	if s.clk.Now().After(t) {
		s.EnqueueAction(ctx, a)
		return
	}
	s.planner.ScheduleAction(ctx, t, a)
}

// moveOverdueActions moves all overdue actions from the planner to the queue.
func (s *SimpleScheduler) moveOverdueActions(ctx context.Context) {
	overdue := s.planner.PopOverdueActions(ctx)

	queue.EnqueueMany(ctx, s.queue, overdue)
}

// RunOne runs one action from the scheduler's queue, returning true if an
// action was run, false if the queue was empty.
func (s *SimpleScheduler) RunOne(ctx context.Context) bool {
	s.moveOverdueActions(ctx)

	if queue.Empty(s.queue) {
		return false
	}

	if a := s.queue.Dequeue(ctx); a != nil {
		events.Run(ctx, a)
		return true
	}
	return false
}

// NextActionTime returns the time of the next action to run, or the current
// time if there are actions to be run in the queue, or util.MaxTime if there
// are no scheduled to run.
func (s *SimpleScheduler) NextActionTime(ctx context.Context) time.Time {
	s.moveOverdueActions(ctx)
	nextScheduled := s.planner.NextActionTime(ctx)

	if !queue.Empty(s.queue) {
		return s.clk.Now()
	}
	return nextScheduled
}
