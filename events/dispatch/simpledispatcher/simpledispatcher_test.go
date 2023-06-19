package simpledispatcher

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p-kad-dht/events"
	ss "github.com/libp2p/go-libp2p-kad-dht/events/scheduler/simplescheduler"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/server/simserver"
	"github.com/stretchr/testify/require"
)

type id string

func (i id) String() string {
	return string(i)
}

func TestSimpleDispatcher(t *testing.T) {
	ctx := context.Background()
	clk := clock.NewMock()

	// creating 5 nodes, with their schedulers
	ids := []address.NodeID{id("a"), id("b"), id("c"), id("d"), id("e")}
	scheds := make(map[address.NodeID]*ss.SimpleScheduler)
	servers := make(map[address.NodeID]simserver.SimServer)
	for _, id := range ids {
		scheds[id] = ss.NewSimpleScheduler(ctx, clk)
		servers[id] = nil
	}

	// creating dispatcher and adding peers
	d := NewSimpleDispatcher(clk)
	for _, id := range ids {
		d.AddPeer(id, scheds[id], servers[id])
	}

	require.Equal(t, servers[ids[0]], d.Server(ids[0]))

	// latencies between nodes
	latencies := [][]int{
		{1},              // b -> a
		{2, 8},           // c -> a, c -> b
		{3, 5, 2},        // d -> a, d -> b, d -> c
		{20, 20, 20, 20}, // e -> a, e -> b, e -> c, e -> d
	}

	// set the latencies
	for i, l := range latencies {
		for j, ll := range l {
			if j > i {
				j += 1
			}
			d.SetLatency(ids[i+1], ids[j], time.Duration(ll)*time.Millisecond)
		}
	}

	// invalid latency set
	d.SetLatency(ids[0], ids[0], 10*time.Millisecond)
	d.SetLatency(ids[1], id("z"), 10*time.Millisecond)

	// check latencies between nodes are correct
	require.Equal(t, time.Duration(0), d.GetLatency(ids[0], ids[0]))
	require.Equal(t, time.Duration(latencies[0][0]*int(time.Millisecond)), d.GetLatency(ids[0], ids[1]))
	require.Equal(t, time.Duration(latencies[0][0]*int(time.Millisecond)), d.GetLatency(ids[1], ids[0]))
	require.Equal(t, time.Duration(latencies[3][0]*int(time.Millisecond)), d.GetLatency(ids[4], ids[0]))
	require.Equal(t, time.Duration(latencies[3][3]*int(time.Millisecond)), d.GetLatency(ids[3], ids[4]))

	// invalid latency (unknown peer)
	require.Equal(t, time.Duration(0), d.GetLatency(ids[0], id("z")))

	// remove peer
	d.RemovePeer(ids[1])
	require.Nil(t, d.Server(ids[1]))

	// add peer again
	d.AddPeer(ids[1], scheds[ids[1]], nil)
	// latency should be reset to 0
	require.Equal(t, time.Duration(0), d.GetLatency(ids[0], ids[1]))

	// create actions to be dispatched
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

	// run one action on each scheduler
	runScheds := func() {
		for _, sched := range scheds {
			sched.RunOne(ctx)
		}
	}

	// dispatch instant action (no latency)
	d.Dispatch(ctx, ids[0], ids[1], actions[0])
	runScheds()
	require.True(t, checks[0])

	d.Dispatch(ctx, ids[0], ids[3], actions[1]) // 3 ms
	d.Dispatch(ctx, ids[2], ids[3], actions[2]) // 2 ms
	d.Dispatch(ctx, ids[4], ids[0], actions[3]) // 20 ms
	d.Dispatch(ctx, ids[2], ids[0], actions[4]) // 2 ms
	d.Dispatch(ctx, ids[4], ids[2], actions[5]) // 20 ms
	d.Dispatch(ctx, ids[4], ids[1], actions[6]) // 0 ms
	d.Dispatch(ctx, ids[3], ids[1], actions[7]) // 0 ms
	clk.Add(4 * time.Millisecond)
	runScheds()

	require.False(t, checks[1]) // actions[2] is prioritary over actions[3] on c
	require.True(t, checks[2])  // actions[2] is prioritary on c
	require.False(t, checks[3]) // actions[4] is prioritary over actions[5] on a
	require.True(t, checks[4])  // actions[4] is prioritary on a
	require.False(t, checks[5]) // it isn't time to run actions[6] on b (20 ms)
	require.True(t, checks[6])  // actions[6] is prioritary on b
	require.False(t, checks[7]) // actions[7] arrives after actions[6] on b

	runScheds()
	require.True(t, checks[1])
	require.False(t, checks[3])
	require.False(t, checks[5])
	require.True(t, checks[7])

	d.DispatchTo(ctx, ids[1], actions[8]) // 0 ms
	runScheds()
	require.True(t, checks[8])
}

func TestDispatchLoop(t *testing.T) {
	ctx := context.Background()
	clk := clock.NewMock()

	// creating 6 nodes, with their schedulers (note that "f" is never used)
	ids := []address.NodeID{id("a"), id("b"), id("c"), id("d"), id("e"), id("f")}
	scheds := make(map[address.NodeID]*ss.SimpleScheduler)
	for _, id := range ids {
		scheds[id] = ss.NewSimpleScheduler(ctx, clk)
	}

	// creating dispatcher and adding peers
	d := NewSimpleDispatcher(clk)
	for _, id := range ids {
		d.AddPeer(id, scheds[id], nil)
	}

	// create actions to be dispatched
	nActions := 10
	actions := make([]events.Action, nActions)

	type checkFormat struct {
		actionId int
		peer     address.NodeID
		time     time.Time
	}

	checks := make([]checkFormat, nActions)

	fnGen := func(i int) func(context.Context) {
		return func(ctx context.Context) {
			a := ctx.Value(ctxActionIdKey).(int)
			p := ctx.Value(ctxPeerKey).(address.NodeID)
			t := ctx.Value(ctxTimeKey).(time.Time)
			checks[i] = checkFormat{a, p, t}
		}
	}
	for i := 0; i < nActions; i++ {
		actions[i] = fnGen(i)
	}

	clk.Set(time.Unix(0, 0))
	f := func(d time.Duration) time.Time {
		return clk.Now().Add(d)
	}
	timings := []time.Time{
		f(3 * time.Millisecond),
		f(10 * time.Millisecond),
		f(100 * time.Millisecond),
		f(100 * time.Millisecond),
		f(100 * time.Millisecond),
		f(500 * time.Millisecond),
		f(500 * time.Millisecond),
		f(10 * time.Second),
		f(time.Minute),
		f(time.Hour),
	}

	d.DispatchDelay(ctx, ids[0], ids[1], actions[1], timings[1]) // b, 10 ms
	d.DispatchDelay(ctx, ids[3], ids[1], actions[0], timings[0]) // b, 3 ms
	d.DispatchDelay(ctx, ids[4], ids[2], actions[2], timings[2]) // c, 100 ms
	d.DispatchDelay(ctx, ids[1], ids[4], actions[3], timings[3]) // e, 100 ms
	d.DispatchDelay(ctx, ids[0], ids[2], actions[4], timings[4]) // c, 100 ms
	d.DispatchDelay(ctx, ids[2], ids[0], actions[5], timings[5]) // a, 500 ms
	d.DispatchDelay(ctx, ids[1], ids[4], actions[6], timings[6]) // e, 500 ms
	d.DispatchDelay(ctx, ids[2], ids[3], actions[7], timings[7]) // d, 10 s
	d.DispatchDelay(ctx, ids[3], ids[0], actions[8], timings[8]) // a, 1 min
	d.DispatchDelay(ctx, ids[4], ids[3], actions[9], timings[9]) // d, 1 hour

	// note that the actions are executed in the order of the actions indexes

	peerOrder := []string{"b", "b", "c", "e", "c", "a", "e", "d", "a", "d"}

	// a: 500ms, 1min
	// b: 3ms, 10ms
	// c: 100ms, 100ms
	// d: 10s, 1h
	// e: 100ms, 500ms

	d.DispatchLoop(ctx)

	for i := 0; i < nActions; i++ {
		require.Equal(t, i, checks[i].actionId)
		require.Equal(t, peerOrder[i], checks[i].peer.String(), i)
		require.Equal(t, timings[i], checks[i].time)
	}
}
