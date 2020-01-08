package uhppote

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	if err := teardown(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	os.Exit(code)
}

func setup() error {
	cmd := exec.Command("docker", "run", "--detach", "--publish", "8000:8000", "--publish", "60000:60000/udp", "--name", "qwerty", "--rm", "integration-tests/simulator")
	out, err := cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("Failed to start Docker simulator instance (%v)", err)
	}

	cmd = exec.Command("docker", "run", "--detach", "--publish", "8081:8080", "--publish", "1883:1883", "--publish", "8883:8883", "--name", "uiop", "--rm", "hivemq/uhppoted")
	out, err = cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("Failed to start Docker HiveMQ instance (%v)", err)
	}

	return nil
}

func teardown() error {
	cmd := exec.Command("docker", "stop", "qwerty")
	out, err := cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("Failed to stop Docker simulator instance (%v)", err)
	}

	cmd = exec.Command("docker", "stop", "uiop")
	out, err = cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("Failed to stop Docker HiveMQ instance (%v)", err)
	}

	return nil
}

func TestMQTTD(t *testing.T) {
	t.Skip("SKIP - not implemented yet")
}
