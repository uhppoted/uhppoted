package mqtt

import (
	"errors"
	"fmt"
	"github.com/uhppoted/uhppoted/src/uhppoted"
)

type ierror struct {
	Err     error  `json:"-"`
	Code    int    `json:"error-code"`
	Message string `json:"message"`
}

var (
	InvalidDeviceID    = ferror(fmt.Errorf("%w: Missing device ID", uhppoted.BadRequest), "Missing device ID")
	InvalidCardNumber  = ferror(fmt.Errorf("%w: Missing/invalid card number", uhppoted.BadRequest), "Missing/invalid card number")
	InvalidDoorID      = ferror(fmt.Errorf("%w: Missing/invalid door ID", uhppoted.BadRequest), "Missing/invalid door ID")
	InvalidDoorDelay   = ferror(fmt.Errorf("%w: Missing/invalid door delay", uhppoted.BadRequest), "Missing/invalid door delay")
	InvalidDoorControl = ferror(fmt.Errorf("%w: Missing/invalid door control", uhppoted.BadRequest), "Missing/invalid door control")
	InvalidEventID     = ferror(fmt.Errorf("%w: Missing/invalid event ID", uhppoted.BadRequest), "Missing/invalid event ID")
	InvalidDateTime    = ferror(fmt.Errorf("%w: Missing/invalid date/time", uhppoted.BadRequest), "Missing/invalid date/time")
)

func (e *ierror) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func ferror(err error, msg string) *ierror {
	status := uhppoted.StatusInternalServerError

	if errors.Is(err, uhppoted.InternalServerError) {
		status = uhppoted.StatusInternalServerError
	} else if errors.Is(err, uhppoted.NotFound) {
		status = uhppoted.StatusNotFound
	}

	return &ierror{
		Err:     err,
		Code:    status,
		Message: msg,
	}
}
