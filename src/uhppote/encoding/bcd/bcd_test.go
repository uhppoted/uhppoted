package bcd

import (
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		s   string
		bcd []byte
	}{
		{"", []byte{}},
		{"1", []byte{0x01}},
		{"12", []byte{0x12}},
		{"123", []byte{0x01, 0x23}},
		{"1234", []byte{0x12, 0x34}},
		{"20191231", []byte{0x20, 0x19, 0x12, 0x31}},
	}

	for _, test := range tests {
		result, err := Encode(test.s)

		if err != nil {
			t.Errorf("bcd.Encode(%s) returned unexpected error: %v", test.s, err)
		}

		if !reflect.DeepEqual(*result, test.bcd) {
			t.Errorf("Invalid packed BCD encoding for '%s': %x", test.s, *result)
		}
	}
}

func TestEncodeInvalidString(t *testing.T) {
	tests := []string{
		" 192837465",
		"192837465 ",
		"1928a37465",
	}

	for _, s := range tests {
		result, err := Encode(s)

		if err == nil {
			t.Errorf("bcd.Encode(%s) should have returned err: %v", s, err)
		}

		if result != nil {
			t.Errorf("bcd.Encode(%s) should have returned 'nil': %v", s, result)
		}
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		s   string
		bcd []byte
	}{
		{"", []byte{}},
		{"01", []byte{0x01}},
		{"12", []byte{0x12}},
		{"0123", []byte{0x01, 0x23}},
		{"1234", []byte{0x12, 0x34}},
		{"20191231", []byte{0x20, 0x19, 0x12, 0x31}},
	}

	for _, test := range tests {
		result, err := Decode(test.bcd)

		if err != nil {
			t.Errorf("bcd.Decode(%x) returned unexpected error: %v", test.bcd, err)
		} else if result != test.s {
			t.Errorf("Invalid BCD decoded string %x: '%s'", test.bcd, result)
		}
	}
}

func TestDecodeInvalidBCD(t *testing.T) {
	tests := [][]byte{
		[]byte{0x0a, 0x12, 0x34},
		[]byte{0x01, 0x2b, 0x34},
		[]byte{0x01, 0x23, 0x3c},
	}

	for _, bytes := range tests {
		result, err := Decode(bytes)

		if err == nil {
			t.Errorf("bcd.Decode(%x) should have returned err: %v", bytes, err)
		}

		if result != "" {
			t.Errorf("bcd.Decode(%x) should have returned 'nil': %v", bytes, result)
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	s := ""

	for i := 0; i < 256; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}

	b.Run("Encode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Encode(s)
		}
	})
}

func BenchmarkDecode(b *testing.B) {
	bytes := make([]byte, 128)

	for i := 0; i < len(bytes); i++ {
		msb := rand.Intn(10)
		lsb := rand.Intn(10)
		bytes[i] = byte(msb<<4 + lsb)
	}

	b.Run("Decode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Decode(bytes)
		}
	})
}
