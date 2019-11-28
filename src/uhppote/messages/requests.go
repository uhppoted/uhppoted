package messages

import (
	"fmt"
	codec "uhppote/encoding/UTO311-L0x"
)

type handler struct {
	factory func() Request
}

var requests = map[byte]func() Request{
	//	0x20: func() Request { return new(GetStatusRequest) },
	//	0x30: func() Request { return new(SetTimeRequest) },
	//	0x32: func() Request { return new(GetTimeRequest) },
	//	0x40: func() Request { return new(OpenDoorRequest) },
	//	0x50: func() Request { return new(PutCardRequest) },
	//	0x52: func() Request { return new(DeleteCardRequest) },
	//	0x54: func() Request { return new(DeleteCardsRequest) },
	//	0x58: func() Request { return new(GetCardsRequest) },
	//	0x5a: func() Request { return new(GetCardByIdRequest) },
	//	0x5c: func() Request { return new(GetCardByIndexRequest) },
	//	0x80: func() Request { return new(SetDoorControlStateRequest) },
	//	0x82: func() Request { return new(GetDoorControlStateRequest) },
	//	0x90: func() Request { return new(SetListenerRequest) },
	//	0x92: func() Request { return new(GetListenerRequest) },
	0x94: func() Request { return new(FindDevicesRequest) },
	//	0x96: func() Request { return new(SetAddressRequest) },
	//	0xb0: func() Request { return new(GetEventRequest) },
	//	0xb2: func() Request { return new(SetEventIndexRequest) },
	//	0xb4: func() Request { return new(GetEventIndexRequest) },
}

var handlers = map[byte]*handler{
	0x20: &handler{
		func() Request { return new(GetStatusRequest) },
	},

	0x30: &handler{
		func() Request { return new(SetTimeRequest) },
	},

	0x32: &handler{
		func() Request { return new(GetTimeRequest) },
	},

	0x40: &handler{
		func() Request { return new(OpenDoorRequest) },
	},

	0x50: &handler{
		func() Request { return new(PutCardRequest) },
	},

	0x52: &handler{
		func() Request { return new(DeleteCardRequest) },
	},

	0x54: &handler{
		func() Request { return new(DeleteCardsRequest) },
	},

	0x58: &handler{
		func() Request { return new(GetCardsRequest) },
	},

	0x5a: &handler{
		func() Request { return new(GetCardByIdRequest) },
	},

	0x5c: &handler{
		func() Request { return new(GetCardByIndexRequest) },
	},

	0x80: &handler{
		func() Request { return new(SetDoorControlStateRequest) },
	},

	0x82: &handler{
		func() Request { return new(GetDoorControlStateRequest) },
	},

	0x90: &handler{
		func() Request { return new(SetListenerRequest) },
	},

	0x92: &handler{
		func() Request { return new(GetListenerRequest) },
	},

	//	0x94: &handler{
	//		func() Request { return new(FindDevicesRequest) },
	//	},

	0x96: &handler{
		func() Request { return new(SetAddressRequest) },
	},

	0xb0: &handler{
		func() Request { return new(GetEventRequest) },
	},

	0xb2: &handler{
		func() Request { return new(SetEventIndexRequest) },
	},

	0xb4: &handler{
		func() Request { return new(GetEventIndexRequest) },
	},
}

func UnmarshalRequest(bytes []byte) (Request, error) {
	if len(bytes) != 64 {
		return nil, fmt.Errorf("Invalid message length %d", len(bytes))
	}

	if bytes[0] != 0x17 {
		return nil, fmt.Errorf("Invalid message type 0x%02x", bytes[0])
	}

	if f := requests[bytes[1]]; f == nil {
		//		return nil, fmt.Errorf("Unknown message type 0x%02x", bytes[1])
	} else {
		response := f()
		err := codec.Unmarshal(bytes, response)
		if err != nil {
			return nil, err
		}

		return response, nil
	}

	if h := handlers[bytes[1]]; h == nil {
		return nil, fmt.Errorf("Unknown message type 0x%02x", bytes[1])
	} else {
		request := h.factory()
		err := codec.Unmarshal(bytes, request)
		if err != nil {
			return nil, err
		}

		return request, nil
	}
}
