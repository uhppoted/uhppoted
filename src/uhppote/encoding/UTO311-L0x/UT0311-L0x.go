package UTO311_L0x

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
	tBool         = reflect.TypeOf(bool(false))
	tByte         = reflect.TypeOf(byte(0))
	tUint16       = reflect.TypeOf(uint16(0))
	tUint32       = reflect.TypeOf(uint32(0))
	tIPv4         = reflect.TypeOf(net.IPv4(0, 0, 0, 0))
	tMAC          = reflect.TypeOf(net.HardwareAddr{})
	tMsgType      = reflect.TypeOf(types.MsgType(0))
	tSerialNumber = reflect.TypeOf(types.SerialNumber(0))
	tVersion      = reflect.TypeOf(types.Version(0))
	tDate         = reflect.TypeOf(types.Date{})
	tDateTime     = reflect.TypeOf(types.DateTime{})
	tSystemDate   = reflect.TypeOf(types.SystemDate{})
	tSystemTime   = reflect.TypeOf(types.SystemTime{})
)

var re = regexp.MustCompile(`offset:\s*([0-9]+)`)
var vre = regexp.MustCompile(`value:\s*0x([0-9a-fA-F]+)`)

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

			// Marshall MsgType field tagged with `uhppote:"value:<value>"`

			if t.Type == tMsgType {
				value := vre.FindStringSubmatch(tag)
				if value == nil {
					panic(errors.New("Cannot marshal message with missing MsgType value"))
				} else {
					v, err := strconv.ParseUint(value[1], 16, 8)
					if err != nil {
						return bytes, err
					}
					bytes[1] = byte(v)
				}

				continue
			}

			matched := re.FindStringSubmatch(tag)

			if matched != nil {
				offset, _ := strconv.Atoi(matched[1])

				switch t.Type {
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

				case tVersion:
					binary.BigEndian.PutUint16(bytes[offset:offset+2], uint16(f.Uint()))

				case tDate:
					slice := reflect.ValueOf(bytes[offset : offset+4])
					f.MethodByName("Encode").Call([]reflect.Value{slice})

				case tDateTime:
					slice := reflect.ValueOf(bytes[offset : offset+7])
					f.MethodByName("Encode").Call([]reflect.Value{slice})

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
		return errors.New(fmt.Sprintf("Invalid start of message - expected 0x17, received 0x%02x", bytes[0]))
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

			if !f.CanSet() {
				continue
			}

			// Unmarshall MsgType field tagged with `uhppote:"value:<value>"`

			if t.Type == tMsgType {
				value := vre.FindStringSubmatch(tag)
				if value == nil {
					panic(errors.New("Cannot unmarshal message with missing MsgType value"))
				} else {
					v, err := strconv.ParseUint(value[1], 16, 8)
					if err != nil {
						return err
					}
					if bytes[1] != byte(v) {
						return errors.New(fmt.Sprintf("Invalid MsgType in message - expected %02x, received 0x%02x", v, bytes[1]))
					}
				}

				f.SetUint(uint64(bytes[1]))
				continue
			}

			// Unmarshall fields tagged with `uhppote:"offset:<offset>"`

			matched := re.FindStringSubmatch(tag)
			if matched == nil {
				continue
			}

			offset, _ := strconv.Atoi(matched[1])

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

			case tDateTime:
				decoded, err := bcd.Decode(bytes[offset : offset+7])
				if err != nil {
					return err
				}

				date, err := time.ParseInLocation("20060102150405", decoded, time.Local)
				if err != nil {
					return err
				}

				f.Field(0).Set(reflect.ValueOf(date))

			case tSystemDate:
				decoded, err := bcd.Decode(bytes[offset : offset+3])
				if err != nil {
					return err
				}

				date, err := time.ParseInLocation("060102", decoded, time.Local)
				if err != nil {
					return err
				}

				f.Field(0).Set(reflect.ValueOf(date))

			case tSystemTime:
				decoded, err := bcd.Decode(bytes[offset : offset+3])
				if err != nil {
					return err
				}

				time, err := time.ParseInLocation("150405", decoded, time.Local)
				if err != nil {
					return err
				}

				f.Field(0).Set(reflect.ValueOf(time))

			default:
				panic(errors.New(fmt.Sprintf("Cannot unmarshal field with type '%v'", t.Type)))
			}
		}
	}

	return nil
}
