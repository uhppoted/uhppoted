package uhppote

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"uhppote/messages"
)

func Marshal(m messages.Message) (*[]byte, error) {
	bytes := make([]byte, 64)

	bytes[0] = 0x17
	bytes[1] = m.Code

	return &bytes, nil
}

func Unmarshal(bytes []byte, m interface{}) error {
	// ... validate message format

	if len(bytes) != 64 {
		return errors.New(fmt.Sprintf("Invalid message length - expected 64 bytes, received %v", len(bytes)))
	}

	if bytes[0] != 0x17 {
		return errors.New(fmt.Sprintf("Invalid start of message - expected 0x17, received 0x%02X", bytes[0]))
	}

	// ... extract fields

	re := regexp.MustCompile(`offset:\s*([0-9]+)`)
	v := reflect.ValueOf(m)
	s := v.Elem()

	if s.Kind() == reflect.Struct {
		N := s.NumField()

		for i := 0; i < N; i++ {
			f := s.Field(i)
			t := s.Type().Field(i)
			tag := t.Tag.Get("uhppote")
			matched := re.FindStringSubmatch(tag)

			if matched != nil && f.CanSet() {
				offset, _ := strconv.Atoi(matched[1])

				switch t.Type {
				case reflect.TypeOf(byte(0)):
					f.SetUint(uint64(bytes[offset]))

				case reflect.TypeOf(uint32(0)):
					f.SetUint(uint64(binary.LittleEndian.Uint32(bytes[offset : offset+4])))

				case reflect.TypeOf(net.IPv4(0, 0, 0, 0)):
					f.SetBytes(net.IPv4(bytes[offset], bytes[offset+1], bytes[offset+2], bytes[offset+3]))

				case reflect.TypeOf(net.HardwareAddr{}):
					f.SetBytes(bytes[offset : offset+6])

				default:
					fmt.Printf("----- MISSING TYPE: %v  %v\n", t.Type.Kind(), t.Type)
				}
			}
		}
	}

	return nil
}
