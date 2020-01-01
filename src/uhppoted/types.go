package uhppoted

import (
	"encoding/json"
	"fmt"
)

type DeviceID uint32

func (id *DeviceID) UnmarshalJSON(bytes []byte) (err error) {
	v := uint32(0)

	if err = json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	if v == 0 {
		err = fmt.Errorf("Invalid DeviceID: %v", v)
		return
	}

	*id = DeviceID(v)

	return
}
