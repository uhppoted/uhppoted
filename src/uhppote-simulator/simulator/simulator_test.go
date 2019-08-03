package simulator

import (
	"net"
	"reflect"
	"testing"
	"time"
	"uhppote-simulator/simulator/entities"
	"uhppote/messages"
	"uhppote/types"
)

func TestHandleGetCardByIndex(t *testing.T) {
	from, _ := types.DateFromString("2019-01-01")
	to, _ := types.DateFromString("2019-12-31")

	request := messages.GetCardByIndexRequest{
		SerialNumber: 12345,
		Index:        2,
	}

	response := messages.GetCardByIndexResponse{
		SerialNumber: 12345,
		CardNumber:   192837465,
		From:         from,
		To:           to,
		Door1:        true,
		Door2:        false,
		Door3:        false,
		Door4:        true,
	}

	testHandle(&request, &response, t)
}

func TestHandleSetDoorDelay(t *testing.T) {
	request := messages.SetDoorDelayRequest{
		SerialNumber: 12345,
		Door:         2,
		Unit:         0x03,
		Delay:        7,
	}

	response := messages.SetDoorDelayResponse{
		SerialNumber: 12345,
		Door:         2,
		Unit:         0x03,
		Delay:        7,
	}

	testHandle(&request, &response, t)
}

func TestHandleGetDoorDelay(t *testing.T) {
	request := messages.GetDoorDelayRequest{
		SerialNumber: 12345,
		Door:         2,
	}

	response := messages.GetDoorDelayResponse{
		SerialNumber: 12345,
		Door:         2,
		Unit:         0x03,
		Delay:        22,
	}

	testHandle(&request, &response, t)
}

func TestHandleSetListener(t *testing.T) {
	request := messages.SetListenerRequest{
		SerialNumber: 12345,
		Address:      net.IPv4(10, 0, 0, 1),
		Port:         43210,
	}

	response := messages.SetListenerResponse{
		SerialNumber: 12345,
		Succeeded:    true,
	}

	testHandle(&request, &response, t)
}

func TestHandleGetListener(t *testing.T) {
	request := messages.GetListenerRequest{
		SerialNumber: 12345,
	}

	response := messages.GetListenerResponse{
		SerialNumber: 12345,
		Address:      net.IPv4(10, 0, 0, 10),
		Port:         43210,
	}

	testHandle(&request, &response, t)
}

// TODO: deferred pending some way to compare Date field
// func TestHandleFindDevices(t *testing.T) {
// 	MAC, _ := net.ParseMAC("00:66:19:39:55:2d")
// 	now := types.Date(time.Now().UTC())
//
// 	request := messages.FindDevicesRequest{}
//
// 	response := messages.FindDevicesResponse{
// 		SerialNumber: 12345,
// 		IpAddress:    net.IPv4(10, 0, 0, 100),
// 		SubnetMask:   net.IPv4(255, 255, 255, 0),
// 		Gateway:      net.IPv4(10, 0, 0, 1),
// 		MacAddress:   types.MacAddress(MAC),
// 		Version:      9876,
// 		Date:         now,
// 	}
//
// 	testHandle(&request, &response, t)
// }

func TestHandleSetAddress(t *testing.T) {
	request := messages.SetAddressRequest{
		SerialNumber: 12345,
		Address:      net.IPv4(10, 0, 0, 100),
		Mask:         net.IPv4(255, 255, 255, 0),
		Gateway:      net.IPv4(10, 0, 0, 1),
		MagicWord:    0x55aaaa55,
	}

	testHandle(&request, nil, t)
}

func TestHandleGetEvent(t *testing.T) {
	timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-08-01 12:34:56", time.Local)

	request := messages.GetEventRequest{
		SerialNumber: 12345,
		Index:        2,
	}

	response := messages.GetEventResponse{
		SerialNumber: 12345,
		Index:        2,
		Type:         0x06,
		Granted:      true,
		Door:         4,
		DoorOpened:   false,
		UserId:       555444321,
		Timestamp:    types.DateTime(timestamp),
		RecordType:   9,
	}

	testHandle(&request, &response, t)
}

func TestHandleSetEventIndex(t *testing.T) {
	request := messages.SetEventIndexRequest{
		SerialNumber: 12345,
		Index:        17,
		MagicWord:    0x55aaaa55,
	}

	response := messages.SetEventIndexResponse{
		SerialNumber: 12345,
		Changed:      true,
	}

	testHandle(&request, &response, t)
}

func TestHandleGetEventIndex(t *testing.T) {
	request := messages.GetEventIndexRequest{
		SerialNumber: 12345,
	}

	response := messages.GetEventIndexResponse{
		SerialNumber: 12345,
		Index:        123,
	}

	testHandle(&request, &response, t)
}

func testHandle(request messages.Request, expected messages.Response, t *testing.T) {
	MAC, _ := net.ParseMAC("00:66:19:39:55:2d")
	from, _ := types.DateFromString("2019-01-01")
	to, _ := types.DateFromString("2019-12-31")
	timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-08-01 12:34:56", time.Local)
	listener := net.UDPAddr{IP: net.IPv4(10, 0, 0, 10), Port: 43210}

	doors := map[uint8]*entities.Door{
		1: &entities.Door{Delay: entities.DelayFromSeconds(11)},
		2: &entities.Door{Delay: entities.DelayFromSeconds(22)},
		3: &entities.Door{Delay: entities.DelayFromSeconds(33)},
		4: &entities.Door{Delay: entities.DelayFromSeconds(44)},
	}

	cards := entities.CardList{
		&entities.Card{100000001, *from, *to, map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
		&entities.Card{192837465, *from, *to, map[uint8]bool{1: true, 2: false, 3: false, 4: true}},
		&entities.Card{100000003, *from, *to, map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
	}

	events := entities.EventList{
		123,
		[]entities.Event{
			entities.Event{
				RecordNumber: 1,
				Type:         0x05,
				Granted:      false,
				Door:         3,
				DoorOpened:   false,
				UserId:       1234567890,
				Timestamp:    types.DateTime(timestamp),
				RecordType:   1,
			},
			entities.Event{
				RecordNumber: 2,
				Type:         0x06,
				Granted:      true,
				Door:         4,
				DoorOpened:   false,
				UserId:       555444321,
				Timestamp:    types.DateTime(timestamp),
				RecordType:   9,
			},
			entities.Event{
				RecordNumber: 3,
				Type:         0x05,
				Granted:      false,
				Door:         3,
				DoorOpened:   false,
				UserId:       1234567890,
				Timestamp:    types.DateTime(timestamp),
				RecordType:   1,
			},
		},
	}

	s := Simulator{
		SerialNumber: 12345,
		IpAddress:    net.IPv4(10, 0, 0, 100),
		SubnetMask:   net.IPv4(255, 255, 255, 0),
		Gateway:      net.IPv4(10, 0, 0, 1),
		MacAddress:   types.MacAddress(MAC),
		Version:      9876,
		Listener:     &listener,
		Cards:        cards,
		Events:       events,
		Doors:        doors,
	}

	response := s.Handle(request)

	if response == nil && expected != nil {
		t.Errorf("Invalid response: Expected: %v, got: %v", expected, response)
		return
	}

	if !reflect.DeepEqual(response, expected) {
		t.Errorf("Incorrect response: Expected: %v, got: %v", expected, response)
	}
}