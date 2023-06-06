package simpleplanner

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p-kad-dht/events"
	"github.com/stretchr/testify/require"
)

func TestSimplePlanner(t *testing.T) {
	ctx := context.Background()
	clk := clock.NewMock()
	p := NewSimplePlanner(clk)

	minTimeStep := time.Nanosecond

	nActions := 10
	actions := make([]events.Action, nActions)
	for i := 0; i < nActions; i++ {
		actions[i] = i
	}

	p.ScheduleAction(ctx, clk.Now().Add(time.Millisecond), actions[0])
	require.Empty(t, p.PopOverdueActions(ctx))

	clk.Add(time.Millisecond)
	require.Empty(t, p.PopOverdueActions(ctx))
	clk.Add(minTimeStep)
	require.Equal(t, actions[:1], p.PopOverdueActions(ctx))
	require.Empty(t, p.PopOverdueActions(ctx))

	p.ScheduleAction(ctx, clk.Now().Add(2*time.Minute), actions[1])
	p.ScheduleAction(ctx, clk.Now().Add(2*time.Second), actions[2])
	p.ScheduleAction(ctx, clk.Now().Add(time.Minute), actions[3])
	p.ScheduleAction(ctx, clk.Now().Add(time.Hour), actions[4])
	require.Empty(t, p.PopOverdueActions(ctx))

	clk.Add(2*time.Second + minTimeStep)
	require.Equal(t, actions[2:3], p.PopOverdueActions(ctx))

	clk.Add(2 * time.Minute)
	require.Equal(t, []events.Action{actions[3], actions[1]}, p.PopOverdueActions(ctx))

	p.ScheduleAction(ctx, clk.Now().Add(time.Second), actions[5])
	clk.Add(time.Second + minTimeStep)
	require.Equal(t, actions[5:6], p.PopOverdueActions(ctx))

	clk.Add(time.Hour)
	require.Equal(t, actions[4:5], p.PopOverdueActions(ctx))

	p.RemoveAction(ctx, actions[0])

	p.ScheduleAction(ctx, clk.Now().Add(time.Second), actions[6])      // 3
	p.ScheduleAction(ctx, clk.Now().Add(time.Microsecond), actions[7]) // 1
	p.ScheduleAction(ctx, clk.Now().Add(time.Hour), actions[8])        // 4
	p.ScheduleAction(ctx, clk.Now().Add(time.Millisecond), actions[9]) // 2

	p.RemoveAction(ctx, actions[9])
	p.RemoveAction(ctx, actions[0])
	clk.Add(time.Second + minTimeStep)

	p.RemoveAction(ctx, actions[6])
	require.Equal(t, actions[7:8], p.PopOverdueActions(ctx))

	p.RemoveAction(ctx, actions[8])
	require.Empty(t, p.PopOverdueActions(ctx))
}
