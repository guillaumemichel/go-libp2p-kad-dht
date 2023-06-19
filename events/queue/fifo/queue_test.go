package fifo

import (
	"context"
	"testing"

	"github.com/libp2p/go-libp2p-kad-dht/events/action"
	"github.com/libp2p/go-libp2p-kad-dht/events/action/testaction"
	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	ctx := context.Background()
	nEvents := 10
	events := make([]action.Action, nEvents)
	for i := 0; i < nEvents; i++ {
		events[i] = testaction.IntAction(i)
	}
	q := NewQueue(ctx)
	if q.Size() != 0 {
		t.Errorf("Expected size 0, got %d", q.Size())
	}
	require.True(t, q.Empty())

	q.Enqueue(ctx, events[0])
	if q.Size() != 1 {
		t.Errorf("Expected size 1, got %d", q.Size())
	}
	require.False(t, q.Empty())

	q.Enqueue(ctx, events[1])
	if q.Size() != 2 {
		t.Errorf("Expected size 2, got %d", q.Size())
	}
	require.False(t, q.Empty())

	newsChan := q.NewsChan()
	require.NotNil(t, newsChan)
	<-newsChan

	e := q.Dequeue(ctx)
	require.Equal(t, e, events[0])
	if q.Size() != 1 {
		t.Errorf("Expected size 1, got %d", q.Size())
	}
	require.False(t, q.Empty())

	e = q.Dequeue(ctx)
	require.Equal(t, e, events[1])
	if q.Size() != 0 {
		t.Errorf("Expected size 0, got %d", q.Size())
	}
	require.True(t, q.Empty())

	e = q.Dequeue(ctx)
	require.Nil(t, e)
	if q.Size() != 0 {
		t.Errorf("Expected size 0, got %d", q.Size())
	}
	require.True(t, q.Empty())

	q.Close()
}
