package kvs

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestKVSLoad(t *testing.T) {
	expected := map[string]string{
		"K1": "A",
		"K2": "B",
		"K3": "C",
	}

	data := `K1  A
K2  B
K3  C`

	kvs := NewKeyValueStore("test", func(v string) (interface{}, error) { return v, nil })
	r := strings.NewReader(data)
	err := kvs.load(r)

	if err != nil {
		t.Fatalf("Unexpected error loading KVS: %v", err)
	}

	for key, v := range expected {
		if value, ok := kvs.Get(key); !ok {
			t.Errorf("%s:%s key-value pair not in store", key, v)
		} else if value != v {
			t.Errorf("Expected %[1]s:%[2]s, got %[1]s:%[3]v", key, v, value)
		}
	}
}

func TestKVSSave(t *testing.T) {
	expected := map[string]string{
		"K1":                    "X",
		"K2":                    "Y",
		"K3":                    "Z",
		"K12345678901234567890": "Q",
	}

	kvs := NewKeyValueStore("test", func(v string) (interface{}, error) { return v, nil })
	kvs.Put("K1", "X")
	kvs.Put("K2", "Y")
	kvs.Put("K3", "Z")
	kvs.Put("K12345678901234567890", "Q")

	buffer := new(bytes.Buffer)
	err := kvs.Save(buffer)

	if err != nil {
		t.Fatalf("Unexpected error saving KVS: %v", err)
	}

	data := buffer.String()

	for k, v := range expected {
		re := regexp.MustCompile(fmt.Sprintf("(?m)^%s\\s{2,}%s$", k, v))
		matched := re.MatchString(data)
		if !matched {
			t.Errorf("%s:%s not matched in saved data", k, v)
		}
	}
}
