package uhppoted

import (
	"errors"
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

var (
	BadRequest          = errors.New("Bad Request")
	NotFound            = errors.New("Not Found")
	InternalServerError = errors.New("INTERNAL SERVER ERROR")
)

type UHPPOTED struct {
	Uhppote         *uhppote.UHPPOTE
	ListenBatchSize int
	Log             *log.Logger
}

func (u *UHPPOTED) debug(tag string, msg interface{}) {
	u.Log.Printf("DEBUG %-12s %v", tag, msg)
}

func (u *UHPPOTED) info(tag string, msg interface{}) {
	u.Log.Printf("INFO  %-12s %v", tag, msg)
}

func (u *UHPPOTED) warn(tag string, err error) {
	u.Log.Printf("WARN  %-12s %v", tag, err)
}
