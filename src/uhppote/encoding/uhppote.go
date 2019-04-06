package uhppote

import (
	"encoding/bcd"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"uhppote/types"
)

var (
	tByte    = reflect.TypeOf(byte(0))
	tUint16  = reflect.TypeOf(uint16(0))
	tUint32  = reflect.TypeOf(uint32(0))
	tIPv4    = reflect.TypeOf(net.IPv4(0, 0, 0, 0))
	tMAC     = reflect.TypeOf(net.HardwareAddr{})
	tVersion = reflect.TypeOf(types.Version(0))
	tDate    = reflect.TypeOf(types.Date{})
)

var re = regexp.MustCompile(`offset:\s*([0-9]+)`)

func Marshal(m interface{}) ([]byte, error) {
	v := reflect.ValueOf(m)

	if v.Type().Kind() == reflect.Ptr {
		return marshal(v.Elem())
	} else {
		return marshal(reflect.Indirect(v))
	}
}

func marshal(s reflect.Value) ([]byte, error) {
	bytes := make([]byte, 64)

	bytes[0] = 0x17

	if s.Kind() == reflect.Struct {
		N := s.NumField()

		for i := 0; i < N; i++ {
			f := s.Field(i)
			t := s.Type().Field(i)
			tag := t.Tag.Get("uhppote")
			matched := re.FindStringSubmatch(tag)

			if matched != nil {
				offset, _ := strconv.Atoi(matched[1])

				switch t.Type {
				case tByte:
					bytes[offset] = byte(f.Uint())

				default:
					panic(errors.New(fmt.Sprintf("Cannot marshal field with type '%v'", t.Type)))
				}
			}
		}
	}

	return bytes, nil
}

func Unmarshal(bytes []byte, m interface{}) error {
	// Validate message format

	if len(bytes) != 64 {
		return errors.New(fmt.Sprintf("Invalid message length - expected 64 bytes, received %v", len(bytes)))
	}

	if bytes[0] != 0x17 {
		return errors.New(fmt.Sprintf("Invalid start of message - expected 0x17, received 0x%02X", bytes[0]))
	}

	// Unmarshall fields tagged with `uhppote:"offset:<offset>"`

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
				case tByte:
					f.SetUint(uint64(bytes[offset]))

				case tUint16:
					f.SetUint(uint64(binary.LittleEndian.Uint16(bytes[offset : offset+2])))

				case tUint32:
					f.SetUint(uint64(binary.LittleEndian.Uint32(bytes[offset : offset+4])))

				case tIPv4:
					f.SetBytes(net.IPv4(bytes[offset], bytes[offset+1], bytes[offset+2], bytes[offset+3]))

				case tMAC:
					f.SetBytes(bytes[offset : offset+6])

				case tVersion:
					f.SetUint(uint64(binary.BigEndian.Uint16(bytes[offset : offset+2])))

				case tDate:
					decoded, err := bcd.Decode(bytes[offset : offset+4])
					if err != nil {
						return err
					}

					date, err := time.ParseInLocation("20060102", decoded, time.Local)
					if err != nil {
						return err
					}

					f.Field(0).Set(reflect.ValueOf(date))

				default:
					panic(errors.New(fmt.Sprintf("Cannot unmarshal field with type '%v'", t.Type)))
				}
			}
		}
	}

	return nil
}
