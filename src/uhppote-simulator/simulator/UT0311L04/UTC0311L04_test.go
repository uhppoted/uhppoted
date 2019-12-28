package UT0311L04

import (
	"net"
	"reflect"
	"testing"
	"time"
	"uhppote-simulator/entities"
	"uhppote/messages"
	"uhppote/types"
)

// TODO: ignore date/time fields
// func TestHandleGetStatus(t *testing.T) {
// 	swipeDateTime, _ := types.DateTimeFromString("2019-08-01 12:34:56")
// 	request := messages.GetStatusRequest{
// 		SerialNumber: 12345,
// 	}
//
// 	response := messages.GetStatusResponse{
// 		SerialNumber:  12345,
// 		LastIndex:     3,
// 		SwipeRecord:   0x00,
// 		Granted:       false,
// 		Door:          3,
// 		DoorOpened:    false,
// 		UserId:        1234567890,
// 		SwipeDateTime: *swipeDateTime,
// 		SwipeReason:   0x05,
// 		Door1State:    false,
// 		Door2State:    false,
// 		Door3State:    false,
// 		Door4State:    false,
// 		Door1Button:   false,
// 		Door2Button:   false,
// 		Door3Button:   false,
// 		Door4Button:   false,
// 		SystemState:   0x00,
// 		//	SystemDate     types.SystemDate   `uhppote:"offset:51"`
// 		//	SystemTime     types.SystemTime   `uhppote:"offset:37"`
// 		PacketNumber:   0,
// 		Backup:         0,
// 		SpecialMessage: 0,
// 		Battery:        0,
// 		FireAlarm:      0,
// 	}
//
// 	testHandle(&request, &response, t)
// }

func TestHandleOpenDoor(t *testing.T) {
	request := messages.OpenDoorRequest{
		SerialNumber: 12345,
		Door:         3,
	}

	response := messages.OpenDoorResponse{
		SerialNumber: 12345,
		Succeeded:    true,
	}

	testHandle(&request, &response, t)
}

func TestHandlePutCardRequest(t *testing.T) {
	from, _ := types.DateFromString("2019-01-01")
	to, _ := types.DateFromString("2019-12-31")
	request := messages.PutCardRequest{
		SerialNumber: 12345,
		CardNumber:   192837465,
		From:         *from,
		To:           *to,
		Door1:        true,
		Door2:        false,
		Door3:        true,
		Door4:        false,
	}

	response := messages.PutCardResponse{
		SerialNumber: 12345,
		Succeeded:    true,
	}

	testHandle(&request, &response, t)
}

func TestHandleDeleteCardRequest(t *testing.T) {
	request := messages.DeleteCardRequest{
		SerialNumber: 12345,
		CardNumber:   192837465,
	}

	response := messages.DeleteCardResponse{
		SerialNumber: 12345,
		Succeeded:    true,
	}

	testHandle(&request, &response, t)
}

func TestHandleDeleteCardsRequest(t *testing.T) {
	request := messages.DeleteCardsRequest{
		SerialNumber: 12345,
		MagicWord:    0x55aaaa55,
	}

	response := messages.DeleteCardsResponse{
		SerialNumber: 12345,
		Succeeded:    true,
	}

	testHandle(&request, &response, t)
}

func TestHandleGetCardsRequest(t *testing.T) {
	request := messages.GetCardsRequest{
		SerialNumber: 12345,
	}

	response := messages.GetCardsResponse{
		SerialNumber: 12345,
		Records:      3,
	}

	testHandle(&request, &response, t)
}

func TestHandleGetCardById(t *testing.T) {
	from, _ := types.DateFromString("2019-01-01")
	to, _ := types.DateFromString("2019-12-31")

	request := messages.GetCardByIdRequest{
		SerialNumber: 12345,
		CardNumber:   192837465,
	}

	response := messages.GetCardByIdResponse{
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

func TestHandleSetDoorControlState(t *testing.T) {
	request := messages.SetDoorControlStateRequest{
		SerialNumber: 12345,
		Door:         2,
		ControlState: 3,
		Delay:        7,
	}

	response := messages.SetDoorControlStateResponse{
		SerialNumber: 12345,
		Door:         2,
		ControlState: 3,
		Delay:        7,
	}

	testHandle(&request, &response, t)
}

func TestHandleGetDoorControlState(t *testing.T) {
	request := messages.GetDoorControlStateRequest{
		SerialNumber: 12345,
		Door:         2,
	}

	response := messages.GetDoorControlStateResponse{
		SerialNumber: 12345,
		Door:         2,
		ControlState: 2,
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
		Result:       9,
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
		1: &entities.Door{ControlState: 3, Delay: entities.DelayFromSeconds(11)},
		2: &entities.Door{ControlState: 2, Delay: entities.DelayFromSeconds(22)},
		3: &entities.Door{ControlState: 3, Delay: entities.DelayFromSeconds(33)},
		4: &entities.Door{ControlState: 3, Delay: entities.DelayFromSeconds(44)},
	}

	cards := entities.CardList{
		&entities.Card{
			CardNumber: 100000001,
			From:       *from,
			To:         *to,
			Doors:      map[uint8]bool{1: false, 2: false, 3: false, 4: false},
		},
		&entities.Card{
			CardNumber: 192837465,
			From:       *from,
			To:         *to,
			Doors:      map[uint8]bool{1: true, 2: false, 3: false, 4: true},
		},
		&entities.Card{
			CardNumber: 100000003,
			From:       *from,
			To:         *to,
			Doors:      map[uint8]bool{1: false, 2: false, 3: false, 4: false},
		},
	}

	events := entities.EventList{
		Index: 123,
		Events: []entities.Event{
			entities.Event{
				RecordNumber: 1,
				Type:         0x05,
				Granted:      false,
				Door:         3,
				DoorOpened:   false,
				UserId:       1234567890,
				Timestamp:    types.DateTime(timestamp),
				Result:       1,
			},
			entities.Event{
				RecordNumber: 2,
				Type:         0x06,
				Granted:      true,
				Door:         4,
				DoorOpened:   false,
				UserId:       555444321,
				Timestamp:    types.DateTime(timestamp),
				Result:       9,
			},
			entities.Event{
				RecordNumber: 3,
				Type:         0x05,
				Granted:      false,
				Door:         3,
				DoorOpened:   false,
				UserId:       1234567890,
				Timestamp:    types.DateTime(timestamp),
				Result:       1,
			},
		},
	}

	txq := make(chan entities.Message, 8)
	src := net.UDPAddr{IP: net.IPv4(10, 0, 0, 1), Port: 12345}

	s := UT0311L04{
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

		txq: txq,
	}

	s.Handle(&src, request)

	if expected != nil {
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(100 * time.Millisecond)
			timeout <- true
		}()

		select {
		case response := <-txq:
			if response.Message == nil {
				t.Errorf("Invalid response: Expected: %v, got: %v", expected, response)
				return
			}

			if !reflect.DeepEqual(response.Message, expected) {
				t.Errorf("Incorrect response: Expected:\n%v, got:s\n%v", expected, response.Message)
			}

		case <-timeout:
			t.Errorf("No response from simulator")
		}
	}
}
