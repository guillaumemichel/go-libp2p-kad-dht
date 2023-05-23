package eventqueue

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	nEvents := 10
	events := make([]*Event, nEvents)
	for i := 0; i < nEvents; i++ {
		events[i] = &Event{}
	}
	q := New()
	if q.Size() != 0 {
		t.Errorf("Expected size 0, got %d", q.Size())
	}
	q.Enqueue(events[0])
	if q.Size() != 1 {
		t.Errorf("Expected size 1, got %d", q.Size())
	}
	q.Enqueue(events[1])
	if q.Size() != 2 {
		t.Errorf("Expected size 2, got %d", q.Size())
	}
	e := q.Dequeue()
	require.Equal(t, e, events[0])
	if q.Size() != 1 {
		t.Errorf("Expected size 1, got %d", q.Size())
	}
	e = q.Dequeue()
	require.Equal(t, e, events[1])
	if q.Size() != 0 {
		t.Errorf("Expected size 0, got %d", q.Size())
	}
	e = q.Dequeue()
	require.Nil(t, e)
	if q.Size() != 0 {
		t.Errorf("Expected size 0, got %d", q.Size())
	}
}
