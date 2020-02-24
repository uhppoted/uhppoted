package types

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
		if ix := IncrementEventIndex(v.index, 100000); ix != v.expected {
			t.Errorf("inccrement %v returned %v, expected %v", v.index, ix, v.expected)
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
		if ix := DecrementEventIndex(v.index, 100000); ix != v.expected {
			t.Errorf("decrement %v returned %v, expected %v", v.index, ix, v.expected)
		}
	}
}
