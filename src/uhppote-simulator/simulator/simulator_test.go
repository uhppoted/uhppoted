package simulator

import (
	"reflect"
	"testing"
	"uhppote-simulator/simulator/entities"
	"uhppote/messages"
	"uhppote/types"
)

func TestHandleGetCardByIndex(t *testing.T) {
	from, _ := types.DateFromString("2019-01-01")
	to, _ := types.DateFromString("2019-12-31")

	expected := messages.GetCardByIndexResponse{
		SerialNumber: 12345,
		CardNumber:   192837465,
		From:         from,
		To:           to,
		Door1:        true,
		Door2:        false,
		Door3:        false,
		Door4:        true,
	}

	request := messages.GetCardByIndexRequest{
		SerialNumber: 12345,
		Index:        2,
	}

	s := setUp()
	response := s.Handle(0x00, &request)

	if response == nil {
		t.Errorf("Invalid response: Expected: %v, got: %v", expected, response)
		return
	}

	if !reflect.DeepEqual(response, &expected) {
		t.Errorf("Incorrect response: Expected: %v, got: %v", expected, response)
	}
}

func TestHandleSetDoorDelay(t *testing.T) {
	expected := messages.SetDoorDelayResponse{
		SerialNumber: 12345,
		Door:         2,
		Unit:         0x03,
		Delay:        7,
	}

	request := messages.SetDoorDelayRequest{
		SerialNumber: 12345,
		Door:         2,
		Unit:         0x03,
		Delay:        7,
	}

	s := setUp()
	response := s.Handle(0x00, &request)

	if response == nil {
		t.Errorf("Invalid response: Expected: %v, got: %v", expected, response)
		return
	}

	if !reflect.DeepEqual(response, &expected) {
		t.Errorf("Incorrect response: Expected: %v, got: %v", expected, response)
	}
}

func TestHandleGetDoorDelay(t *testing.T) {
	expected := messages.GetDoorDelayResponse{
		SerialNumber: 12345,
		Door:         2,
		Unit:         0x03,
		Delay:        22,
	}

	request := messages.GetDoorDelayRequest{
		SerialNumber: 12345,
		Door:         2,
	}

	s := setUp()
	response := s.Handle(0x00, &request)

	if response == nil {
		t.Errorf("Invalid response: Expected: %v, got: %v", expected, response)
		return
	}

	if !reflect.DeepEqual(response, &expected) {
		t.Errorf("Incorrect response: Expected: %v, got: %v", expected, response)
	}
}

func setUp() Simulator {
	from, _ := types.DateFromString("2019-01-01")
	to, _ := types.DateFromString("2019-12-31")

	cards := entities.CardList{
		&entities.Card{100000001, *from, *to, map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
		&entities.Card{192837465, *from, *to, map[uint8]bool{1: true, 2: false, 3: false, 4: true}},
		&entities.Card{100000003, *from, *to, map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
	}

	doors := map[uint8]*entities.Door{
		1: &entities.Door{Delay: entities.DelayFromSeconds(11)},
		2: &entities.Door{Delay: entities.DelayFromSeconds(22)},
		3: &entities.Door{Delay: entities.DelayFromSeconds(33)},
		4: &entities.Door{Delay: entities.DelayFromSeconds(44)},
	}

	return Simulator{
		SerialNumber: 12345,
		Cards:        cards,
		Doors:        doors,
	}
}
