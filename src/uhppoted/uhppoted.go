package uhppoted

import (
	"log"
	"net/http"
	"uhppote"
)

const (
	StatusOK                  = http.StatusOK
	StatusBadRequest          = http.StatusBadRequest
	StatusNotFound            = http.StatusNotFound
	StatusInternalServerError = http.StatusInternalServerError
)

type UHPPOTED struct {
	Uhppote *uhppote.UHPPOTE
	Log     *log.Logger
}

func (u *UHPPOTED) log(tag string, deviceID uint32, msg string) {
	u.Log.Printf("%-5s %-12d %s", tag, deviceID, msg)
}

func (u *UHPPOTED) debug(operation string, msg interface{}) {
	u.Log.Printf("DEBUG %-20s %v", operation, msg)
}

func (u *UHPPOTED) info(deviceID uint32, operation string, msg interface{}) {
	u.Log.Printf("INFO   %-12d %-20s %v", deviceID, operation, msg)
}

func (u *UHPPOTED) warn(deviceID uint32, operation string, err error) {
	u.Log.Printf("WARN  %-12d %-20s %v", deviceID, operation, err)
}
