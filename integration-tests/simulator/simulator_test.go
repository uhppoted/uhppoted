package uhppote

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"
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
	cmd := exec.Command("docker", "run", "--detach", "--publish", "8000:8000", "--publish", "60000:60000/udp", "--name", "simulator", "--rm", "integration-tests/simulator")
	out, err := cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("Failed to start Docker simulator instance (%v)", err)
	}

	return nil
}

func teardown() error {
	cmd := exec.Command("docker", "stop", "simulator")
	out, err := cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("Failed to stop Docker simulator instance (%v)", err)
	}

	return nil
}

//	curl -X POST "http://127.0.0.1:8000/uhppote/simulator/405419896/swipe" -H "accept: application/json" -H "Content-Type: application/json" -d "{\"door\":3,\"card-number\":65538}"
func TestSwipe(t *testing.T) {
	url := "http://127.0.0.1:8000/uhppote/simulator/405419896/swipe"
	payload := strings.NewReader(`{"door":3, "card-number":65538}`)

	rq, err := http.NewRequest("POST", url, payload)
	if err != nil {
		t.Fatalf("Error creating swipe request (%v)", err)
	}

	rq.Header.Add("cache-control", "no-cache")
	rq.Header.Add("accept", "application/json")
	rq.Header.Add("content-type", "application/json")

	response, err := http.DefaultClient.Do(rq)
	if err != nil {
		t.Fatalf("Error POST'ing swipe request (%v)", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Incorrect response status - expected:%v, got: %v", http.StatusOK, response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Error reading swipe response (%v)", err)
	}

	expected := struct {
		AccessGranted bool   `json:"access-granted"`
		DoorOpened    bool   `json:"door-opened"`
		Message       string `json:"message"`
	}{
		AccessGranted: false,
		DoorOpened:    false,
		Message:       "Access denied",
	}

	reply := struct {
		AccessGranted bool   `json:"access-granted"`
		DoorOpened    bool   `json:"door-opened"`
		Message       string `json:"message"`
	}{}

	if err := json.Unmarshal(body, &reply); err != nil {
		t.Fatalf("Error parsing response (%v)", err)
	}

	if !reflect.DeepEqual(reply, expected) {
		t.Errorf("Incorrect reply to 'swipe' request\n   expected: %+v\n   got:      %+v", expected, reply)
	}
}
