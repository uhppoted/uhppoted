package uhppoted

import (
	"context"
	"log"
	"net/http"
)

const (
	StatusOK                  = http.StatusOK
	StatusBadRequest          = http.StatusBadRequest
	StatusNotFound            = http.StatusNotFound
	StatusInternalServerError = http.StatusInternalServerError
)

type Service interface {
	Send(ctx context.Context, message interface{})
	Reply(ctx context.Context, response interface{})
	Oops(ctx context.Context, operation string, message string, errorCode int)
}

type Request interface {
	DeviceID() (*uint32, error)
	DeviceDoor() (*uint32, *uint8, error)
	DeviceDoorDelay() (*uint32, *uint8, *uint8, error)
	DeviceDoorControl() (*uint32, *uint8, *string, error)
}

type UHPPOTED struct {
	Service Service
}

type Device struct {
	ID uint32 `json:"id"`
}

func (u *UHPPOTED) log(ctx context.Context, tag string, deviceID uint32, msg string) {
	ctx.Value("log").(*log.Logger).Printf("%-5s %-12d %s\n", tag, deviceID, msg)
}

func (u *UHPPOTED) debug(ctx context.Context, operation string, msg interface{}) {
	ctx.Value("log").(*log.Logger).Printf("DEBUG %-20s %v\n", operation, msg)
}

func (u *UHPPOTED) info(ctx context.Context, deviceID uint32, operation string, rq Request) {
	ctx.Value("log").(*log.Logger).Printf("INFO   %-12d %-20s %s\n", deviceID, operation, rq)
}

func (u *UHPPOTED) warn(ctx context.Context, deviceID uint32, operation string, err error) {
	ctx.Value("log").(*log.Logger).Printf("WARN  %-12d %-20s %v\n", deviceID, operation, err)
}

func (u *UHPPOTED) send(ctx context.Context, message interface{}) {
	u.Service.Send(ctx, message)
}

func (u *UHPPOTED) reply(ctx context.Context, response interface{}) {
	u.Service.Reply(ctx, response)
}

func (u *UHPPOTED) oops(ctx context.Context, operation string, message string, errorCode int) {
	u.Service.Oops(ctx, operation, message, errorCode)
}
