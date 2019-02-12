package bcd

import (
	"errors"
	"reflect"
	"testing"
)

var tests = []struct {
	s        string
	expected []byte
	err      error
}{
	{"", []byte{}, nil},
	{"1", []byte{0x01}, nil},
	{"12", []byte{0x12}, nil},
	{"123", []byte{0x01, 0x23}, nil},
	{"1234", []byte{0x12, 0x34}, nil},
	{"192837465", []byte{0x01, 0x92, 0x83, 0x74, 0x65}, nil},
	{" 192837465", nil, errors.New("Invalid numeric string ' 192837465'")},
	{"192837465 ", nil, errors.New("Invalid numeric string '192837465 '")},
	{"1928a37465", nil, errors.New("Invalid numeric string '1928a37465'")},
}

func TestEncode(t *testing.T) {
	for _, test := range tests {
		result, err := Encode(test.s)

		if err != nil && !reflect.DeepEqual(err, test.err) {
			t.Errorf("bcd.Encode(%s) returned incorrect error: %v", test.s, err)
		}

		if test.expected == nil && result != nil {
			t.Errorf("Invalid packed BCD encoding for '%s': %x", test.s, *result)
		} else if test.expected != nil && !reflect.DeepEqual(*result, test.expected) {
			t.Errorf("Invalid packed BCD encoding for '%s': %x", test.s, *result)
		}
	}
}
