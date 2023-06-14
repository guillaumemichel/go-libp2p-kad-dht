package simpleplanner

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p-kad-dht/events"
	"github.com/libp2p/go-libp2p-kad-dht/internal/util"
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

	a0 := p.ScheduleAction(ctx, clk.Now().Add(time.Millisecond), actions[0])
	require.Empty(t, p.PopOverdueActions(ctx))

	clk.Add(time.Millisecond)
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

	p.RemoveAction(ctx, a0)

	a6 := p.ScheduleAction(ctx, clk.Now().Add(time.Second), actions[6])      // 3
	p.ScheduleAction(ctx, clk.Now().Add(time.Microsecond), actions[7])       // 1
	a8 := p.ScheduleAction(ctx, clk.Now().Add(time.Hour), actions[8])        // 4
	a9 := p.ScheduleAction(ctx, clk.Now().Add(time.Millisecond), actions[9]) // 2

	p.RemoveAction(ctx, a9)
	p.RemoveAction(ctx, a0)
	clk.Add(time.Second + minTimeStep)

	p.RemoveAction(ctx, a6)
	require.Equal(t, actions[7:8], p.PopOverdueActions(ctx))

	p.RemoveAction(ctx, a8)
	require.Empty(t, p.PopOverdueActions(ctx))
}

func TestNextActionTime(t *testing.T) {
	ctx := context.Background()
	clk := clock.NewMock()
	p := NewSimplePlanner(clk)

	clk.Set(time.Unix(0, 0))

	ti := p.NextActionTime(ctx)
	require.Equal(t, util.MaxTime, ti)

	t0 := clk.Now().Add(time.Second)
	p.ScheduleAction(ctx, t0, 0)
	ti = p.NextActionTime(ctx)
	require.Equal(t, t0, ti)

	t1 := clk.Now().Add(time.Hour)
	p.ScheduleAction(ctx, t1, 1)
	ti = p.NextActionTime(ctx)
	require.Equal(t, t0, ti)

	t2 := clk.Now().Add(time.Millisecond)
	p.ScheduleAction(ctx, t2, 2)
	ti = p.NextActionTime(ctx)
	require.Equal(t, t2, ti)

	require.Equal(t, 0, len(p.PopOverdueActions(ctx)))

	clk.Add(time.Millisecond)
	ti = p.NextActionTime(ctx)
	require.Equal(t, t2, ti)

	require.Equal(t, 1, len(p.PopOverdueActions(ctx)))
	ti = p.NextActionTime(ctx)
	require.Equal(t, t0, ti)

	clk.Add(time.Hour)
	require.Equal(t, 2, len(p.PopOverdueActions(ctx)))
	ti = p.NextActionTime(ctx)
	require.Equal(t, util.MaxTime, ti)
}
