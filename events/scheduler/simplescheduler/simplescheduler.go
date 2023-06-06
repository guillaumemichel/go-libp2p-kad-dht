package simplescheduler

import (
	"context"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p-kad-dht/events"
	"github.com/libp2p/go-libp2p-kad-dht/events/planner"
	"github.com/libp2p/go-libp2p-kad-dht/events/planner/simpleplanner"
	"github.com/libp2p/go-libp2p-kad-dht/events/queue"
	"github.com/libp2p/go-libp2p-kad-dht/events/queue/chanqueue"
)

type SimpleScheduler struct {
	clk clock.Clock

	queue   queue.EventQueue
	planner planner.ActionPlanner
}

func NewSimpleScheduler(ctx context.Context, clk clock.Clock) *SimpleScheduler {
	return &SimpleScheduler{
		clk: clk,

		queue:   chanqueue.NewChanQueue(ctx, 100),
		planner: simpleplanner.NewSimplePlanner(clk),
	}
}

func (s *SimpleScheduler) Now() time.Time {
	return s.clk.Now()
}

func (s *SimpleScheduler) EnqueueAction(ctx context.Context, a events.Action) {
	s.queue.Enqueue(a)
}

func (s *SimpleScheduler) ScheduleAction(ctx context.Context, t time.Time, a events.Action) {
	if s.clk.Now().After(t) {
		s.EnqueueAction(ctx, a)
		return
	}
	s.planner.ScheduleAction(ctx, t, a)
}

func (s *SimpleScheduler) moveOverdueActions(ctx context.Context) {
	overdue := s.planner.PopOverdueActions(ctx)

	queue.EnqueueMany(s.queue, overdue)
}

func (s *SimpleScheduler) RunOne(ctx context.Context) {
	s.moveOverdueActions(ctx)

	if queue.Empty(s.queue) {
		return
	}

	if a := s.queue.Dequeue(); a != nil {
		events.Run(ctx, a)
	}
}
