package messages

import (
	"errors"
	"fmt"
	codec "uhppote/encoding/UTO311-L0x"
)

type handler struct {
	factory func() Request
}

var handlers = map[byte]*handler{
	0x20: &handler{
		func() Request { return new(GetStatusRequest) },
	},

	//0x30: &handler{
	//	func() Request { return new(uhppote.SetTimeRequest) },
	//},

	//0x32: &handler{
	//	func() Request { return new(uhppote.GetTimeRequest) },
	//},

	//0x40: &handler{
	//	func() Request { return new(uhppote.OpenDoorRequest) },
	//},

	//0x50: &handler{
	//	func() Request { return new(uhppote.PutCardRequest) },
	//},

	//0x52: &handler{
	//	func() Request { return new(uhppote.DeleteCardRequest) },
	//},

	//0x54: &handler{
	//	func() Request { return new(uhppote.DeleteCardsRequest) },
	//},

	//0x58: &handler{
	//	func() Request { return new(uhppote.GetCardsRequest) },
	//},

	//0x5a: &handler{
	//	func() Request { return new(uhppote.GetCardByIdRequest) },
	//},

	//0x5c: &handler{
	//	func() Request { return new(uhppote.GetCardByIndexRequest) },
	//},

	//0x80: &handler{
	//	func() Request { return new(uhppote.SetDoorDelayRequest) },
	//},

	//0x82: &handler{
	//	func() Request { return new(uhppote.GetDoorDelayRequest) },
	//},

	//0x94: &handler{
	//	func() Request { return new(uhppote.FindDevicesRequest) },
	//},

	//0x96: &handler{
	//	func() Request { return new(uhppote.SetAddressRequest) },
	//},

	//0xb0: &handler{
	//	func() Request { return new(uhppote.GetEventRequest) },
	//},

	0xb2: &handler{
		func() Request { return new(SetEventIndexRequest) },
	},

	0xb4: &handler{
		func() Request { return new(GetEventIndexRequest) },
	},
}

func UnmarshalRequest(bytes []byte) (*Request, error) {
	if len(bytes) != 64 {
		return nil, errors.New(fmt.Sprintf("Invalid message length %d", len(bytes)))
	}

	if bytes[0] != 0x17 {
		return nil, errors.New(fmt.Sprintf("Invalid message type 0x%02x", bytes[0]))
	}

	if h := handlers[bytes[1]]; h == nil {
		return nil, errors.New(fmt.Sprintf("Unknown message type 0x%02x", bytes[1]))
	} else {
		request := h.factory()
		err := codec.Unmarshal(bytes, request)
		if err != nil {
			return nil, err
		}

		return &request, nil
	}
}
