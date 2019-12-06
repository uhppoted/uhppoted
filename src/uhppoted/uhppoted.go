package uhppoted

import (
	"context"
	"log"
	"net/http"
	"time"
	"uhppote/types"
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
	DateTime() (*time.Time, error)
	DeviceDoor() (*uint32, *uint8, error)
	DeviceDoorDelay() (*uint32, *uint8, *uint8, error)
	DeviceDoorControl() (*uint32, *uint8, *string, error)
	DeviceCardID() (*uint32, *uint32, error)
	DeviceCard() (*uint32, *types.Card, error)
}

type UHPPOTED struct {
	Service Service
}

func (u *UHPPOTED) log(ctx context.Context, tag string, deviceId uint32, msg string) {
	ctx.Value("log").(*log.Logger).Printf("%-5s %-12d %s\n", tag, deviceId, msg)
}

func (u *UHPPOTED) debug(ctx context.Context, deviceId uint32, operation string, rq interface{}) {
	ctx.Value("log").(*log.Logger).Printf("DEBUG %-20s %-12d %v\n", operation, deviceId, rq)
}

func (u *UHPPOTED) info(ctx context.Context, deviceId uint32, operation string, rq Request) {
	ctx.Value("log").(*log.Logger).Printf("INFO   %-12d %-20s %s\n", deviceId, operation, rq)
}

func (u *UHPPOTED) warn(ctx context.Context, deviceId uint32, operation string, err error) {
	ctx.Value("log").(*log.Logger).Printf("WARN  %-12d %-20s %v\n", deviceId, operation, err)
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
