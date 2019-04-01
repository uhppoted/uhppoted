package uhppote

import (
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) Search() ([]types.Device, error) {
	devices := []types.Device{}
	request, err := messages.NewFindDevicesRequest()

	if err != nil {
		return nil, err
	}

	reply, err := u.Exec(request)

	if err != nil {
		return nil, err
	}

	result, err := messages.NewSearch(reply)

	if err != nil {
		return nil, err
	}

	devices = append(devices, result.Device)

	return devices, nil
}
