package UTO311_L0x

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"uhppote/types"
)

type Marshaler interface {
	MarshalUT0311L0x() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalUT0311L0x([]byte) error
}

var (
	tBool         = reflect.TypeOf(bool(false))
	tByte         = reflect.TypeOf(byte(0))
	tUint16       = reflect.TypeOf(uint16(0))
	tUint32       = reflect.TypeOf(uint32(0))
	tIPv4         = reflect.TypeOf(net.IPv4(0, 0, 0, 0))
	tMAC          = reflect.TypeOf(net.HardwareAddr{})
	tMsgType      = reflect.TypeOf(types.MsgType(0))
	tSerialNumber = reflect.TypeOf(types.SerialNumber(0))
	tDatePtr      = reflect.TypeOf((*types.Date)(nil))
)

var re = regexp.MustCompile(`offset:\s*([0-9]+)`)
var vre = regexp.MustCompile(`value:\s*(0[xX])?([0-9a-fA-F]+)`)

func Marshal(m interface{}) ([]byte, error) {
	v := reflect.ValueOf(m)

	if v.Type().Kind() == reflect.Ptr {
		return marshal(v.Elem())
	} else {
		return marshal(reflect.Indirect(v))
	}
}

func marshal(s reflect.Value) ([]byte, error) {
	msgType, err := getMsgType(s)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Cannot marshal message: %v", err)))
	}

	bytes := make([]byte, 64)
	bytes[0] = 0x17
	bytes[1] = msgType

	if s.Kind() == reflect.Struct {
		N := s.NumField()

		for i := 0; i < N; i++ {
			f := s.Field(i)
			t := s.Type().Field(i)
			tag := t.Tag.Get("uhppote")

			matched := re.FindStringSubmatch(tag)

			if matched != nil {
				offset, _ := strconv.Atoi(matched[1])

				// Use MarshalUT0311L0x() if implemented.
				if m, ok := f.Interface().(Marshaler); ok {
					// If f is a pointer type and the value is nil skips this field, leaving the buffer 'as is'
					if f.Kind() != reflect.Ptr || !f.IsNil() {
						if b, err := m.MarshalUT0311L0x(); err == nil {
							copy(bytes[offset:offset+len(b)], b)
						}
					}

					continue
				}

				// Marshal built-in types
				switch t.Type {
				case tMsgType:
				case tByte:
					value := vre.FindStringSubmatch(tag)
					if value != nil {
						v, err := strconv.ParseUint(value[1], 16, 8)
						if err != nil {
							return bytes, err
						}
						bytes[offset] = byte(v)
					} else {
						bytes[offset] = byte(f.Uint())
					}

				case tUint16:
					binary.LittleEndian.PutUint16(bytes[offset:offset+4], uint16(f.Uint()))

				case tUint32:
					binary.LittleEndian.PutUint32(bytes[offset:offset+4], uint32(f.Uint()))

				case tBool:
					if f.Bool() {
						bytes[offset] = 0x01
					} else {
						bytes[offset] = 0x00
					}

				case tIPv4:
					copy(bytes[offset:offset+4], f.MethodByName("To4").Call([]reflect.Value{})[0].Bytes())

				case tMAC:
					copy(bytes[offset:offset+6], f.Bytes())

				case tSerialNumber:
					binary.LittleEndian.PutUint32(bytes[offset:offset+4], uint32(f.Uint()))

				default:
					panic(errors.New(fmt.Sprintf("Cannot marshal field with type '%v'", t.Type)))
				}
			}
		}
	}

	return bytes, nil
}

func getMsgType(s reflect.Value) (byte, error) {
	if s.Kind() == reflect.Struct {
		N := s.NumField()

		for i := 0; i < N; i++ {
			t := s.Type().Field(i)
			tag := t.Tag.Get("uhppote")

			if t.Type == tMsgType {
				value := vre.FindStringSubmatch(tag)

				if value == nil {
					return 0x00, errors.New(fmt.Sprintf("MsgType field has invalid tag:<%v> - expected 'value:<value>'", tag))
				}

				if value[1] == "0x" || value[1] == "0X" {
					v, err := strconv.ParseUint(value[2], 16, 8)
					if err != nil {
						return 0x00, errors.New(fmt.Sprintf("Invalid MsgType value: %v", err))
					}

					return byte(v), nil
				}

				v, err := strconv.ParseUint(value[2], 10, 8)
				if err != nil {
					return 0x00, errors.New(fmt.Sprintf("Invalid MsgType value: %v", err))
				}

				return byte(v), nil
			}
		}
	}

	return 0x00, errors.New("Missing MsgType field")
}

func Unmarshal(bytes []byte, m interface{}) error {
	// Validate message format

	if len(bytes) != 64 {
		return errors.New(fmt.Sprintf("Invalid message length - expected 64 bytes, received %v", len(bytes)))
	}

	if bytes[0] != 0x17 {
		return errors.New(fmt.Sprintf("Invalid start of message - expected 0x17, received 0x%02x", bytes[0]))
	}

	// Unmarshall fields tagged with `uhppote:"..."`

	v := reflect.ValueOf(m)
	s := v.Elem()

	msgType, err := getMsgType(s)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Cannot unmarshal message: %v", err)))
	}

	if bytes[1] != msgType {
		return errors.New(fmt.Sprintf("Invalid MsgType in message - expected %02x, received 0x%02x", msgType, bytes[1]))
	}

	if s.Kind() == reflect.Struct {
		N := s.NumField()

		for i := 0; i < N; i++ {
			f := s.Field(i)
			t := s.Type().Field(i)
			tag := t.Tag.Get("uhppote")

			if !f.CanSet() {
				continue
			}

			// Unmarshall MsgType field

			if t.Type == tMsgType {
				f.SetUint(uint64(msgType))
				continue
			}

			// Unmarshall fields tagged with `uhppote:"offset:<offset>"`

			matched := re.FindStringSubmatch(tag)
			if matched == nil {
				continue
			}

			offset, _ := strconv.Atoi(matched[1])

			// Use UnmarshalUT0311L0x() if implemented
			if u, ok := f.Addr().Interface().(Unmarshaler); ok {
				if err := u.UnmarshalUT0311L0x(bytes[offset:]); err == nil {
					continue
				}
			}

			// Unmarshal built-in types
			switch t.Type {
			case tBool:
				if bytes[offset] == 0x01 {
					f.SetBool(true)
				} else if bytes[offset] == 0x00 {
					f.SetBool(false)
				} else {
					return errors.New(fmt.Sprintf("Invalid boolean value in message: %02x:", bytes[offset]))
				}

			case tByte:
				value := vre.FindStringSubmatch(tag)
				if value != nil {
					v, err := strconv.ParseUint(value[1], 16, 8)
					if err != nil {
						return err
					}
					if bytes[offset] != byte(v) {
						return errors.New(fmt.Sprintf("Invalid value in message - expected %02x, received 0x%02x", v, bytes[offset]))
					}
				}

				f.SetUint(uint64(bytes[offset]))

			case tUint16:
				f.SetUint(uint64(binary.LittleEndian.Uint16(bytes[offset : offset+2])))

			case tUint32:
				f.SetUint(uint64(binary.LittleEndian.Uint32(bytes[offset : offset+4])))

			case tIPv4:
				f.SetBytes(net.IPv4(bytes[offset], bytes[offset+1], bytes[offset+2], bytes[offset+3]))

			case tMAC:
				f.SetBytes(bytes[offset : offset+6])

			case tSerialNumber:
				f.SetUint(uint64(binary.LittleEndian.Uint32(bytes[offset : offset+4])))

			case tDatePtr:
				d := reflect.New(reflect.TypeOf(types.Date{}))
				if u, ok := d.Interface().(Unmarshaler); ok {
					if err := u.UnmarshalUT0311L0x(bytes[offset:]); err == nil {
						f.Set(d)
					}
				} else {
					panic(errors.New(fmt.Sprintf("Cannot unmarshal field with type '%v'", t.Type)))
				}

			default:
				panic(errors.New(fmt.Sprintf("Cannot unmarshal field with type '%v'", t.Type)))
			}
		}
	}

	return nil
}
