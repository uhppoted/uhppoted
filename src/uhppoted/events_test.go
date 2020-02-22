package uhppoted

import (
	"testing"
)

func TestEventRollover(t *testing.T) {
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
		if ix := decrement(v.index, 100000); ix != v.expected {
			t.Errorf("decrement %v returned %v, expected %v", v.index, ix, v.expected)
		}
	}
}
