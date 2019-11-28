package uhppoted

import (
	"context"
	"uhppote"
)

type DoorDelay struct {
	Device struct {
		ID    uint32 `json:"id"`
		Door  uint8  `json:"door"`
		Delay uint8  `json:"delay"`
	} `json:"device"`
}

type DoorControl struct {
	Device struct {
		ID           uint32 `json:"id"`
		Door         uint8  `json:"door"`
		ControlState string `json:"control"`
	} `json:"device"`
}

func (u *UHPPOTED) GetDoorDelay(ctx context.Context, rq Request) {
	u.debug(ctx, 0, "get-door-delay", rq)

	id, door, err := rq.DeviceDoor()
	if err != nil {
		u.warn(ctx, 0, "get-door-delay", err)
		u.oops(ctx, "get-door-delay", "Error retrieving door delay (invalid device/door)", StatusBadRequest)
		return
	}

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(*id, *door)
	if err != nil {
		u.warn(ctx, *id, "get-door-delay", err)
		u.oops(ctx, "get-door-delay", "Error retrieving door delay", StatusInternalServerError)
		return
	}

	response := DoorDelay{
		struct {
			ID    uint32 `json:"id"`
			Door  uint8  `json:"door"`
			Delay uint8  `json:"delay"`
		}{
			ID:    *id,
			Door:  result.Door,
			Delay: result.Delay,
		},
	}

	u.reply(ctx, response)
}

func (u *UHPPOTED) SetDoorDelay(ctx context.Context, rq Request) {
	u.debug(ctx, 0, "set-door-delay", rq)

	id, err := rq.DeviceId()
	if err != nil {
		u.warn(ctx, 0, "set-door-delay", err)
		u.oops(ctx, "get-door-delay", "Error setting door delay (invalid device ID)", StatusBadRequest)
		return
	}

	door, delay, err := rq.DoorDelay()
	if err != nil {
		u.warn(ctx, *id, "set-door-delay", err)
		u.oops(ctx, "set-door-delay", "Error setting door delay (invalid door/delay)", StatusBadRequest)
		return
	}

	state, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(*id, *door)
	if err != nil {
		u.warn(ctx, *id, "set-door-delay", err)
		u.oops(ctx, "set-door-delay", "Error setting door delay (invalid delay)", StatusInternalServerError)
		return
	}

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).SetDoorControlState(*id, *door, state.ControlState, *delay)
	if err != nil {
		u.warn(ctx, *id, "set-door-delay", err)
		u.oops(ctx, "set-door-delay", "Error setting door delay", StatusInternalServerError)
		return
	}

	response := DoorDelay{
		struct {
			ID    uint32 `json:"id"`
			Door  uint8  `json:"door"`
			Delay uint8  `json:"delay"`
		}{
			ID:    *id,
			Door:  result.Door,
			Delay: result.Delay,
		},
	}

	u.reply(ctx, response)
}

func (u *UHPPOTED) GetDoorControl(ctx context.Context, rq Request) {
	u.debug(ctx, 0, "get-door-control", rq)

	id, door, err := rq.DeviceDoor()
	if err != nil {
		u.warn(ctx, *id, "get-door-control", err)
		u.oops(ctx, "get-door-control", "Error retrieving door control (invalid device/door)", StatusBadRequest)
		return
	}

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(*id, *door)
	if err != nil {
		u.warn(ctx, *id, "get-door-control", err)
		u.oops(ctx, "get-door-control", "Error retrieving door control", StatusInternalServerError)
		return
	}

	lookup := map[uint8]string{
		1: "normally open",
		2: "normally closed",
		3: "controlled",
	}

	response := DoorControl{
		struct {
			ID           uint32 `json:"id"`
			Door         uint8  `json:"door"`
			ControlState string `json:"control"`
		}{
			ID:           *id,
			Door:         result.Door,
			ControlState: lookup[result.ControlState],
		},
	}

	u.reply(ctx, response)
}

// func setDoorControl(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	states := map[string]uint8{
// 		"normally open":   1,
// 		"normally closed": 2,
// 		"controlled":      3,
// 	}
//
// 	lookup := map[uint8]string{
// 		1: "normally open",
// 		2: "normally closed",
// 		3: "controlled",
// 	}
//
// 	deviceId := ctx.Value("device-id").(uint32)
// 	door := ctx.Value("door").(uint8)
//
// 	blob, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		warn(ctx, deviceId, "set-door-control", err)
// 		http.Error(w, "Error reading request", http.StatusInternalServerError)
// 		return
// 	}
//
// 	body := struct {
// 		ControlState string `json:"control"`
// 	}{}
//
// 	err = json.Unmarshal(blob, &body)
// 	if err != nil {
// 		warn(ctx, deviceId, "set-door-control", err)
// 		http.Error(w, "Invalid request format", http.StatusBadRequest)
// 		return
// 	} else if _, ok := states[body.ControlState]; !ok {
// 		warn(ctx, deviceId, "set-door-control", fmt.Errorf("Invalid request control value: '%s'", body.ControlState))
// 		http.Error(w, fmt.Sprintf("Invalid request value: '%s'", body.ControlState), http.StatusBadRequest)
// 		return
// 	}
//
// 	state, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(deviceId, door)
// 	if err != nil {
// 		warn(ctx, deviceId, "set-door-control", err)
// 		http.Error(w, "Error setting door control", http.StatusInternalServerError)
// 		return
// 	}
//
// 	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).SetDoorControlState(deviceId, door, states[body.ControlState], state.Delay)
// 	if err != nil {
// 		warn(ctx, deviceId, "set-door-control", err)
// 		http.Error(w, "Error setting door control", http.StatusInternalServerError)
// 		return
// 	}
//
// 	response := struct {
// 		ControlState string `json:"control"`
// 	}{
// 		ControlState: lookup[result.ControlState],
// 	}
//
// 	reply(ctx, w, response)
// }
