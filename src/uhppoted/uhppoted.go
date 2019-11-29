package uhppoted

import (
	"context"
	"log"
	"net/http"
	"time"
)

const (
	StatusBadRequest          = http.StatusBadRequest
	StatusNotFound            = http.StatusNotFound
	StatusInternalServerError = http.StatusInternalServerError
)

type Service interface {
	Reply(ctx context.Context, response interface{})
	Oops(ctx context.Context, operation string, message string, errorCode int)
}

type Request interface {
	DeviceId() (*uint32, error)
	DateTime() (*time.Time, error)
	DeviceDoor() (*uint32, *uint8, error)
	DeviceDoorDelay() (*uint32, *uint8, *uint8, error)
	DeviceDoorControl() (*uint32, *uint8, *string, error)
	DeviceCard() (*uint32, *uint32, error)
}

type UHPPOTED struct {
	Service Service
}

func (u *UHPPOTED) debug(ctx context.Context, deviceId int, operation string, rq Request) {
	ctx.Value("log").(*log.Logger).Printf("DEBUG %-12d %-20s %s\n", deviceId, operation, rq)
}

func (u *UHPPOTED) warn(ctx context.Context, deviceId uint32, operation string, err error) {
	ctx.Value("log").(*log.Logger).Printf("WARN  %-12d %-20s %v\n", deviceId, operation, err)
}

func (u *UHPPOTED) reply(ctx context.Context, response interface{}) {
	u.Service.Reply(ctx, response)
}

func (u *UHPPOTED) oops(ctx context.Context, operation string, message string, errorCode int) {
	u.Service.Oops(ctx, operation, message, errorCode)
}
