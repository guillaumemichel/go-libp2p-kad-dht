package simplescheduler

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p-kad-dht/events"
	"github.com/libp2p/go-libp2p-kad-dht/events/scheduler"
	"github.com/libp2p/go-libp2p-kad-dht/internal/util"
	"github.com/stretchr/testify/require"
)

func TestSimpleScheduler(t *testing.T) {
	ctx := context.Background()
	clk := clock.NewMock()

	sched := NewSimpleScheduler(ctx, clk)

	require.Equal(t, clk.Now(), sched.Now())

	nActions := 10
	actions := make([]events.Action, nActions)
	checks := make([]bool, nActions)

	fnGen := func(i int) func(context.Context) {
		return func(ctx context.Context) {
			checks[i] = true
		}
	}

	for i := 0; i < nActions; i++ {
		actions[i] = fnGen(i)
	}

	sched.EnqueueAction(ctx, actions[0])
	require.False(t, checks[0])
	sched.RunOne(ctx)
	require.True(t, checks[0])

	scheduler.ScheduleActionIn(ctx, sched, time.Second, actions[1])
	require.False(t, checks[1])
	sched.EnqueueAction(ctx, actions[2])
	clk.Add(2 * time.Second)

	sched.RunOne(ctx)
	require.True(t, checks[2])
	require.False(t, checks[1])
	sched.RunOne(ctx)
	require.True(t, checks[1])
	sched.RunOne(ctx)

	scheduler.ScheduleActionIn(ctx, sched, -1*time.Second, actions[3])
	require.False(t, checks[3])
	sched.RunOne(ctx)
	require.True(t, checks[3])

	sched.ScheduleAction(ctx, clk.Now().Add(-1*time.Nanosecond), actions[4])
	require.False(t, checks[4])
	sched.RunOne(ctx)
	require.True(t, checks[4])

	sched.ScheduleAction(ctx, clk.Now().Add(time.Second), actions[5])
	sched.RunOne(ctx)
	require.False(t, checks[5])
	clk.Add(time.Second)
	require.Equal(t, clk.Now(), sched.NextActionTime(ctx))
	sched.RunOne(ctx)
	require.True(t, checks[5])

	t6 := clk.Now().Add(time.Second)
	a6 := sched.ScheduleAction(ctx, t6, actions[6])
	require.Equal(t, t6, sched.NextActionTime(ctx))
	sched.RemovePlannedAction(ctx, a6)
	clk.Add(time.Second)
	sched.RunOne(ctx)
	require.False(t, checks[6])
	// empty queue
	require.Equal(t, util.MaxTime, sched.NextActionTime(ctx))

}
