package uhppote

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

func setup() {
}

func teardown() {
}

func TestLoadTSV(t *testing.T) {
	t.Skip("SKIP - not implemented yet")
}
