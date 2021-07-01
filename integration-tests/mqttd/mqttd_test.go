package uhppote

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

type object map[string]interface{}

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
	cmd := exec.Command("docker", "run", "--detach", "--publish", "8081:8080", "--publish", "1883:1883", "--publish", "8883:8883", "--name", "hivemq", "--rm", "integration-tests/hivemq")
	out, err := cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("Failed to start Docker HiveMQ instance (%v)", err)
	}

	cmd = exec.Command("docker", "run", "--detach", "--publish", "8000:8000", "--publish", "60000:60000/udp", "--name", "simulator", "--rm", "integration-tests/simulator")
	out, err = cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("Failed to start Docker simulator instance (%v)", err)
	}

	time.Sleep(15 * time.Second)

	cmd = exec.Command("docker", "run", "--detach", "--name", "mqttd", "--rm", "integration-tests/mqttd")
	out, err = cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("Failed to start Docker HiveMQ instance (%v)", err)
	}

	return nil
}

func teardown() error {
	cmd := exec.Command("docker", "stop", "mqttd")
	out, err := cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("Failed to stop Docker MQTTD instance (%v)", err)
	}

	cmd = exec.Command("docker", "stop", "simulator")
	out, err = cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("Failed to stop Docker simulator instance (%v)", err)
	}

	cmd = exec.Command("docker", "stop", "hivemq")
	out, err = cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("Failed to stop Docker HiveMQ instance (%v)", err)
	}

	return nil
}

func send(topic, msg string) error {
	cmd := exec.Command("mqtt", "publish", "--topic", topic, "--message", msg)
	out, err := cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("Failed to publish to MQTT topic (%v)", err)
	}

	return nil
}

func listen(queue chan object) error {
	cmd := exec.Command("mqtt", "subscribe", "--topic", "uhppoted/reply/#")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Failed to subscribe to MQTT topic (%v)", err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("Failed to subscribe to MQTT topic (%v)", err)
	}

	o := map[string]interface{}{}
	if err := json.NewDecoder(stdout).Decode(&o); err != nil {
		return fmt.Errorf("Error subscribing to MQTT topic (%v)", err)
	}

	queue <- o

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("Error subscribing to MQTT topic (%v)", err)
	}

	return nil
}

func TestGetDevices(t *testing.T) {
	timeout := make(chan string, 1)
	queue := make(chan object, 1)
	topic := "uhppoted/gateway/requests/devices:get"
	request := `{ "message": { "request": { "request-id":"AH173635G3", "client-id":"QWERTY54", "reply-to":"uhppoted/reply/97531" }}}`

	time.Sleep(10 * time.Second)
	go listen(queue)
	go func() {
		time.Sleep(5 * time.Second)
		timeout <- "timeout"
	}()

	send(topic, request)

	select {
	case <-timeout:
		t.Errorf("TIMEOUT")

	case reply := <-queue:
		fmt.Printf("REPLY: %v", reply)
	}
}
