package uhppoted

import (
	"testing"
)

func TestIncrementEventIndex(t *testing.T) {
	vector := []struct {
		index    uint32
		expected uint32
	}{
		{0, 1},
		{1, 2},
		{19, 20},
		{99999, 100000},
		{100000, 1},
		{100001, 1},
	}

	for _, v := range vector {
		ix := EventIndex(v.index)
		jx := ix.increment(100000)

		if uint32(ix) != v.index {
			t.Errorf("increment %v updated index %v, expected %v", v.index, ix, v.index)
		}

		if uint32(jx) != v.expected {
			t.Errorf("increment %v returned %v, expected %v", v.index, jx, v.expected)
		}
	}
}

func TestDecrementEventIndex(t *testing.T) {
	vector := []struct {
		index    uint32
		expected uint32
	}{
		{100000, 99999},
		{19, 18},
		{1, 100000},
		{0, 100000},
	}

	for _, v := range vector {
		ix := EventIndex(v.index)
		jx := ix.decrement(100000)

		if uint32(ix) != v.index {
			t.Errorf("decrement %v updated %v, expected %v", v.index, ix, v.index)
		}

		if uint32(jx) != v.expected {
			t.Errorf("decrement %v returned %v, expected %v", v.index, jx, v.expected)
		}
	}
}
