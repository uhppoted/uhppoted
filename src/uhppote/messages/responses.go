package messages

import (
	"fmt"
	codec "uhppote/encoding/UTO311-L0x"
)

var responses = map[byte]func() Response{
	0x20: func() Response { return new(GetStatusResponse) },
	0x30: func() Response { return new(SetTimeResponse) },
	0x32: func() Response { return new(GetTimeResponse) },
	0x40: func() Response { return new(OpenDoorResponse) },
	0x50: func() Response { return new(PutCardResponse) },
	0x52: func() Response { return new(DeleteCardResponse) },
	0x54: func() Response { return new(DeleteCardsResponse) },
	0x58: func() Response { return new(GetCardsResponse) },
	0x5a: func() Response { return new(GetCardByIdResponse) },
	0x5c: func() Response { return new(GetCardByIndexResponse) },
	0x80: func() Response { return new(SetDoorControlStateResponse) },
	//	0x82: func() Response { return new(GetDoorControlStateResponse) },
	//	0x90: func() Response { return new(SetListenerResponse) },
	//	0x92: func() Response { return new(GetListenerResponse) },
	0x94: func() Response { return new(FindDevicesResponse) },
	//	0x96: func() Response { return new(SetAddressResponse) },
	//	0xb0: func() Response { return new(GetEventResponse) },
	//	0xb2: func() Response { return new(SetEventIndexResponse) },
	//	0xb4: func() Response { return new(GetEventIndexResponse) },
}

func UnmarshalResponse(bytes []byte) (Response, error) {
	if len(bytes) != 64 {
		return nil, fmt.Errorf("Invalid message length %d", len(bytes))
	}

	if bytes[0] != 0x17 {
		return nil, fmt.Errorf("Invalid protocol ID 0x%02x", bytes[0])
	}

	f := responses[bytes[1]]
	if f == nil {
		return nil, fmt.Errorf("Unknown message type 0x%02x", bytes[1])
	}

	response := f()
	if err := codec.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}
